package payroll

import (
	"context"
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	// Simulate returns the per-session evidence + rolled totals for
	// the given employee/date-range. Read-only.
	Simulate(
		ctx context.Context,
		adminID string,
		startDate, endDate time.Time,
	) (*entity.SimulatePayrollResultDto, error)

	// Save persists the header + components atomically.
	Save(ctx context.Context, dto *entity.EmployeeSalaryDto) (*entity.EmployeeSalaryDto, error)

	// FindAll paginates saved payroll runs.
	FindAll(ctx context.Context, req *entity.EmployeeSalaryFindAllRequest) (*pagination.ResultPagination, error)

	// FindOne returns a single payroll run with its components.
	FindOne(ctx context.Context, id string) (*entity.EmployeeSalaryDto, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Simulate(
	ctx context.Context,
	adminID string,
	startDate, endDate time.Time,
) (*entity.SimulatePayrollResultDto, error) {
	var sessions []model.StockSession
	if err := r.db.WithContext(ctx).
		Where("employee_id = ? AND date >= ? AND date <= ?", adminID, startDate, endDate).
		Order("date ASC").
		Find(&sessions).Error; err != nil {
		return nil, err
	}

	// Resolve the employee's company + the salary_component rows
	// once for this run. The per-session breakdown is then derived
	// live from the SAME rule that RecomputeSalary would have
	// computed at close time, so a historic session whose
	// persisted snapshot is stale (e.g. salary_component was added
	// later) still reports the right number on a fresh simulate.
	//
	// Failure paths (no company binding / no salary_component
	// rows) silently fall back to zero — same as the runtime close
	// path — so the simulate numbers always agree with what a
	// fresh close would have produced.
	var companyID string
	if err := r.db.WithContext(ctx).
		Model(&model.AdminCompany{}).
		Select("company_id").
		Where("admin_id = ?", adminID).
		Scan(&companyID).Error; err != nil {
		companyID = ""
	}
	var components []entity.SalaryComponentDto
	if companyID != "" {
		var rows []model.SalaryComponent
		if err := r.db.WithContext(ctx).
			Where("company_id = ?", companyID).
			Find(&rows).Error; err == nil {
			for i := range rows {
				components = append(components, *entity.NewSalaryComponentDtoFromModel(&rows[i]))
			}
		}
	}

	// Pull every cash_debt row for this employee in the same date
	// range, indexed by YYYY-MM-DD so each session can be paired
	// with the advances that landed on its day. Skipping org-scoping
	// here is intentional — the simulator is admin-scoped, the
	// cash_debt rows already enforce org binding at write time.
	var debtRows []model.CashDebt
	if err := r.db.WithContext(ctx).
		Where("admin_id_employee = ? AND date >= ? AND date <= ?", adminID, startDate, endDate).
		Order("date ASC").
		Find(&debtRows).Error; err != nil {
		debtRows = nil
	}
	cashDebtByDate := make(map[string][]model.CashDebt, len(debtRows))
	for _, d := range debtRows {
		key := d.Date.Format("2006-01-02")
		cashDebtByDate[key] = append(cashDebtByDate[key], d)
	}
	toDebtDto := func(d model.CashDebt) entity.SimulatePayrollSessionCashDebtDto {
		return entity.SimulatePayrollSessionCashDebtDto{
			ID:            d.ID,
			Date:          d.Date.Format("2006-01-02"),
			Amount:        d.Amount,
			PaymentMethod: d.PaymentMethod,
			Notes:         d.Notes,
		}
	}

	out := &entity.SimulatePayrollResultDto{
		AdminIDEmployee: adminID,
		StartDate:       startDate.Format("2006-01-02"),
		EndDate:         endDate.Format("2006-01-02"),
		Sessions:        make([]entity.SimulatePayrollSessionDto, 0, len(sessions)),
		SessionCount:    len(sessions),
	}
	for _, ss := range sessions {
		// Live-recompute salary for this session — same inputs
		// RecomputeSalary(components, TotalItems) would have used
		// at close time. Falls back to the persisted snapshot when
		// no salary config exists, matching the runtime close
		// contract exactly.
		var meal, attendance, bonus, totalSalary float64
		if len(components) > 0 {
			dto := &entity.StockSessionDto{
				TotalItems: ss.TotalItems,
			}
			dto.RecomputeSalary(components, ss.TotalItems)
			meal = dto.MealAllowance
			attendance = dto.Attendance
			bonus = dto.BonusTarget
			totalSalary = dto.TotalSalary
		} else {
			meal = ss.MealAllowance
			attendance = ss.Attendance
			bonus = ss.BonusTarget
			totalSalary = ss.TotalSalary
		}
		// Attach every cash_debt row whose date matches this
		// session's date. A session that happens on a date with no
		// advances still gets an empty (non-nil) slice so the
		// frontend can render an empty state without null-checks.
		sessionDate := ss.Date.Format("2006-01-02")
		debtsForSession := cashDebtByDate[sessionDate]
		debtDtos := make([]entity.SimulatePayrollSessionCashDebtDto, 0, len(debtsForSession))
		var sessionDebtTotal float64
		for _, d := range debtsForSession {
			debtDtos = append(debtDtos, toDebtDto(d))
			sessionDebtTotal += d.Amount
		}

		out.Sessions = append(out.Sessions, entity.SimulatePayrollSessionDto{
			SessionID:     ss.ID,
			Date:          sessionDate,
			Status:        ss.Status,
			TotalSales:    ss.TotalSales,
			Commission:    ss.TotalCommission,
			MealAllowance: meal,
			BonusTarget:   bonus,
			TotalSalary:   totalSalary,
			CashDebts:     debtDtos,
		})
		out.TotalCommission += ss.TotalCommission
		out.TotalMealAllowance += meal
		out.TotalBonusTarget += bonus
		out.TotalAttendance += attendance
		// If the same cash_debt row falls on a date that has two
		// sessions (rare: half-day split), it would be counted twice
		// here. We accept the duplication as the conservative choice
		// so the operator can see the row attached to BOTH closes.
		out.TotalCashDebt += sessionDebtTotal
	}
	out.TotalSalary = out.TotalMealAllowance + out.TotalAttendance +
		out.TotalCommission + out.TotalBonusTarget
	// Cash_debt is money the driver already owes the company on
	// the same range — net it out alongside cash_receipt so
	// remaining reflects what the company still owes the driver.
	out.RemainingSalary = out.TotalSalary - out.TotalCashReceipt - out.TotalCashDebt
	return out, nil
}

func (r *repository) Save(
	ctx context.Context,
	dto *entity.EmployeeSalaryDto,
) (*entity.EmployeeSalaryDto, error) {
	return r.saveTx(ctx, dto)
}

// saveTx wraps the header + components in a transaction so the
// employee_salary row only lands if every component row also lands.
func (r *repository) saveTx(
	ctx context.Context,
	dto *entity.EmployeeSalaryDto,
) (*entity.EmployeeSalaryDto, error) {
	var out *entity.EmployeeSalaryDto
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		header := dto.ToModel()
		if err := tx.Create(header).Error; err != nil {
			return err
		}
		components := make([]model.EmployeeSalaryComponent, 0, len(dto.Components))
		for _, c := range dto.Components {
			m := c.ToModel()
			m.EmployeeSalaryID = header.ID
			components = append(components, *m)
		}
		if len(components) > 0 {
			if err := tx.Create(&components).Error; err != nil {
				return err
			}
		}
		out = entity.NewEmployeeSalaryDtoFromModel(header)
		out.Components = make([]entity.EmployeeSalaryComponentDto, 0, len(components))
		for i := range components {
			out.Components = append(out.Components, *entity.NewEmployeeSalaryComponentDtoFromModel(&components[i]))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (r *repository) FindOne(
	ctx context.Context,
	id string,
) (*entity.EmployeeSalaryDto, error) {
	var m model.EmployeeSalary
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	out := entity.NewEmployeeSalaryDtoFromModel(&m)
	var cs []model.EmployeeSalaryComponent
	if err := r.db.WithContext(ctx).Where("employee_salary_id = ?", id).Find(&cs).Error; err != nil {
		return nil, err
	}
	for i := range cs {
		out.Components = append(out.Components, *entity.NewEmployeeSalaryComponentDtoFromModel(&cs[i]))
	}
	return out, nil
}

func (r *repository) FindAll(
	ctx context.Context,
	req *entity.EmployeeSalaryFindAllRequest,
) (*pagination.ResultPagination, error) {
	var rows []model.EmployeeSalary = make([]model.EmployeeSalary, 0)
	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		q := r.db.WithContext(ctx).Model(&model.EmployeeSalary{})
		if req.AdminIDEmployee != "" {
			q = q.Where("admin_id_employee = ?", req.AdminIDEmployee)
		}
		return q
	}, &pagination.TableRequest{
		Request:       req,
		QueryField:    []string{},
		Data:          &rows,
		AllowedFields: []string{"start_date", "end_date", "created_at"},
	})
	if err != nil {
		return nil, err
	}
	result := dataTable.(*pagination.ResultPagination)
	hits := result.Data.(*[]model.EmployeeSalary)
	out := make([]*entity.EmployeeSalaryDto, 0, len(*hits))
	for i := range *hits {
		out = append(out, entity.NewEmployeeSalaryDtoFromModel(&(*hits)[i]))
	}
	return &pagination.ResultPagination{
		Data:        out,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}
