package stocksession

import (
	"context"
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, dto *entity.StockSessionDto) (*entity.StockSessionDto, error)
	Get(ctx context.Context, id string) (*entity.StockSessionDto, error)
	GetByEmployeeDate(ctx context.Context, employeeID, date string) (*entity.StockSessionDto, error)
	Update(ctx context.Context, dto *entity.StockSessionDto) (*entity.StockSessionDto, error)
	FindAll(ctx context.Context, req *entity.StockSessionFindAllRequest) (*pagination.ResultPagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, dto *entity.StockSessionDto) (*entity.StockSessionDto, error) {
	var result *entity.StockSessionDto
	err := r.db.Transaction(func(tx *gorm.DB) error {
		m := dto.ToModel()

		// Two-phase write so we never rely on GORM association
		// auto-insert. Insert the parent first (without items), then
		// the child stock_session_item rows explicitly. This:
		//   1. Makes the wire-level intent clear in Postgres logs.
		//   2. Avoids accidental item creation when the request
		//      carries a partial payload (e.g. an empty `items`
		//      slice).
		//   3. Lets us reject orphan item rows cleanly without GORM
		//      silent cascade.
		parent := &model.StockSession{
			OrganizationID: m.OrganizationID,
			Date:           m.Date,
			EmployeeID:     m.EmployeeID,
			Status:         m.Status,
			OpenedAt:       m.OpenedAt,
			ClosedAt:       m.ClosedAt,
			TotalSales:     m.TotalSales,
			TotalCash:      m.TotalCash,
			TotalQris:      m.TotalQris,
			TotalOther:     m.TotalOther,
			TotalPayment:   m.TotalPayment,
			Difference:     m.Difference,
			TotalItems:     m.TotalItems,
			// Salary columns: open-time values are 0 today because
			// the close path is what triggers RecomputeSalary, but
			// we deliberately round-trip them so a future change to
			// Open-time resolution (e.g. seed the row with the
			// employee's company default) lands in the DB without a
			// schema migration.
			TotalCommission: m.TotalCommission,
			MealAllowance:   m.MealAllowance,
			Attendance:      m.Attendance,
			BonusTarget:     m.BonusTarget,
			TotalSalary:     m.TotalSalary,
			// cash_debt is 0 at open time; close path writes it.
			CashDebt:  m.CashDebt,
			Notes:     m.Notes,
			CreatedBy: m.CreatedBy,
		}
		if err := tx.Omit("Items", "Payments", "Adjustments", "Logs").Create(parent).Error; err != nil {
			return err
		}

		// Insert stock_session_item rows explicitly. We deliberately
		// do NOT rely on GORM's association auto-save here — that
		// would also try to INSERT into the item table whenever
		// `dto.Items[i].Item` is populated. The Omit("Item") is
		// belt-and-braces so even a stray nested Item cannot leak.
		// All we persist is the item_id reference + the row's own
		// snapshot fields.
		for _, it := range dto.Items {
			item := it.ToModel()
			item.SessionID = parent.ID
			if err := tx.Omit("Item").Create(item).Error; err != nil {
				return err
			}
		}

		// Reload with associations so the call-site gets a fully
		// hydrated DTO back (items, payments, etc.).
		var reloaded model.StockSession
		if err := tx.
			Preload("Employee").
			Preload("Items").
			Preload("Items.Item").
			Where("id = ?", parent.ID).First(&reloaded).Error; err != nil {
			return err
		}
		result = entity.NewStockSessionDtoFromModel(&reloaded)
		return nil
	})
	return result, err
}

func (r *repository) Get(ctx context.Context, id string) (*entity.StockSessionDto, error) {
	var m *model.StockSession
	if err := r.db.
		Preload("Employee").
		Preload("Items").
		Preload("Items.Item").
		Preload("Payments").
		Preload("Adjustments").
		Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return entity.NewStockSessionDtoFromModel(m), nil
}

func (r *repository) GetByEmployeeDate(ctx context.Context, employeeID, date string) (*entity.StockSessionDto, error) {
	var m *model.StockSession
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}
	if err := r.db.
		Preload("Employee").
		Preload("Items").
		Preload("Items.Item").
		Where("employee_id = ? AND date = ?", employeeID, parsedDate).First(&m).Error; err != nil {
		return nil, err
	}
	return entity.NewStockSessionDtoFromModel(m), nil
}

func (r *repository) Update(ctx context.Context, dto *entity.StockSessionDto) (*entity.StockSessionDto, error) {
	return r.updateInTx(dto)
}

func (r *repository) updateInTx(dto *entity.StockSessionDto) (*entity.StockSessionDto, error) {
	var result *entity.StockSessionDto
	err := r.db.Transaction(func(tx *gorm.DB) error {
		m := dto.ToModel()
		// Replace children atomically
		if err := tx.Where("session_id = ?", m.ID).Delete(&model.StockSessionItem{}).Error; err != nil {
			return err
		}
		if err := tx.Where("session_id = ?", m.ID).Delete(&model.PaymentDetail{}).Error; err != nil {
			return err
		}
		if err := tx.Where("session_id = ?", m.ID).Delete(&model.CashAdjustment{}).Error; err != nil {
			return err
		}
		if err := tx.Save(m).Error; err != nil {
			return err
		}
		// Reload with associations
		var reloaded model.StockSession
		if err := tx.
			Preload("Employee").
			Preload("Items").
			Preload("Items.Item").
			Preload("Payments").
			Preload("Adjustments").
			Where("id = ?", m.ID).First(&reloaded).Error; err != nil {
			return err
		}
		result = entity.NewStockSessionDtoFromModel(&reloaded)
		return nil
	})
	return result, err
}

func (r *repository) FindAll(ctx context.Context, req *entity.StockSessionFindAllRequest) (*pagination.ResultPagination, error) {
	var m []model.StockSession = make([]model.StockSession, 0)
	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		q := r.db.Model(&model.StockSession{})
		if req.FindAllRequest.OrganizationData.ID != "" {
			q = q.Where("organization_id = ?", req.FindAllRequest.OrganizationData.ID)
		}
		return q
	}, &pagination.TableRequest{
		Request:       req,
		QueryField:    []string{},
		Data:          &m,
		AllowedFields: []string{"date", "status", "total_sales"},
	})
	if err != nil {
		return nil, err
	}
	result := dataTable.(*pagination.ResultPagination)
	results := result.Data.(*[]model.StockSession)
	data := make([]*entity.StockSessionDto, 0)
	for i := range *results {
		// shallow hydrate (no nested) for list view — keeps payload small
		row := entity.NewStockSessionDtoFromModel(&(*results)[i])
		// Preload items via separate query (lightweight count only)
		data = append(data, row)
	}
	return &pagination.ResultPagination{
		Data:        data,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}
