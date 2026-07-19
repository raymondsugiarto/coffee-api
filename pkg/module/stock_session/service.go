package stocksession

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
	"gorm.io/gorm"
)

// Service is the business logic for daily stock sessions of coffee-cart drivers.
type Service interface {
	Open(ctx context.Context, dto *entity.StockSessionDto, actorID string) (*entity.StockSessionDto, error)
	Get(ctx context.Context, id string) (*entity.StockSessionDto, error)
	GetByEmployeeDate(ctx context.Context, employeeID, date string) (*entity.StockSessionDto, error)
	Update(ctx context.Context, dto *entity.StockSessionDto, actorID string) (*entity.StockSessionDto, error)
	Close(ctx context.Context, id string, dto *entity.StockSessionDto, actorID string) (*entity.StockSessionDto, error)
	FindAll(ctx context.Context, req *entity.StockSessionFindAllRequest) (*pagination.ResultPagination, error)
	GetDashboard(ctx context.Context) (*entity.DashboardSummaryDto, error)
	GetDailyReport(ctx context.Context, date string) (*entity.DailyReportDto, error)
	GetMonthlyReport(ctx context.Context, year int, month int) (*entity.MonthlyReportDto, error)
	GetTopProducts(ctx context.Context, from, to string, limit int) ([]entity.TopProductRowDto, error)
	GetEmployeePerformance(ctx context.Context, from, to string) ([]entity.EmployeePerformanceRowDto, error)
}

type service struct {
	repo Repository
	db   *gorm.DB
}

func NewService(repo Repository, db *gorm.DB) Service {
	return &service{repo: repo, db: db}
}

func (s *service) Open(ctx context.Context, dto *entity.StockSessionDto, actorID string) (*entity.StockSessionDto, error) {
	if dto.EmployeeID == "" {
		return nil, status.New(status.BadRequest, errors.New("employeeId is required"))
	}
	if dto.Date == "" {
		return nil, status.New(status.BadRequest, errors.New("date is required"))
	}
	if len(dto.Items) == 0 {
		return nil, status.New(status.BadRequest, errors.New("at least one item is required"))
	}
	// Open-time items carry a morning OutQty (>= 1, validated on the
	// wire DTO). The internal DTO has already been normalised from
	// the open wire shape by the handler, so OutQty is always set here.

	// Reject duplicate (employee, date)
	if existing, _ := s.repo.GetByEmployeeDate(ctx, dto.EmployeeID, dto.Date); existing != nil {
		return nil, status.New(status.BadRequest, errors.New("stock session already exists for this employee and date"))
	}

	// Verify driver is EMPLOYEE
	var driver *model.Admin
	if err := s.db.Where("id = ? AND admin_type = ?", dto.EmployeeID, "EMPLOYEE").First(&driver).Error; err != nil {
		return nil, status.New(status.BadRequest, errors.New("driver not found or not an employee"))
	}

	// Hydrate item snapshots
	if err := s.hydrateItemSnapshots(ctx, dto); err != nil {
		return nil, err
	}

	dto.OrganizationID = shared.GetOrganization(ctx).ID
	dto.Status = entity.StockSessionStatusOpen
	dto.OpenedAt = time.Now()
	dto.CreatedBy = actorID
	dto.RecomputeTotals()
	s.resolveAndApplySalary(ctx, dto)

	result, err := s.repo.Create(ctx, dto)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) Get(ctx context.Context, id string) (*entity.StockSessionDto, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) GetByEmployeeDate(ctx context.Context, employeeID, date string) (*entity.StockSessionDto, error) {
	return s.repo.GetByEmployeeDate(ctx, employeeID, date)
}

func (s *service) Update(ctx context.Context, dto *entity.StockSessionDto, actorID string) (*entity.StockSessionDto, error) {
	existing, err := s.repo.Get(ctx, dto.ID)
	if err != nil {
		return nil, err
	}
	if existing.Status == entity.StockSessionStatusClosed {
		return nil, status.New(status.BadRequest, errors.New("cannot update closed session"))
	}
	if err := s.hydrateItemSnapshots(ctx, dto); err != nil {
		return nil, err
	}
	dto.OrganizationID = existing.OrganizationID
	dto.Status = existing.Status
	dto.OpenedAt = existing.OpenedAt
	dto.CreatedBy = existing.CreatedBy
	dto.RecomputeTotals()
	s.resolveAndApplySalary(ctx, dto)

	result, err := s.repo.Update(ctx, dto)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) Close(ctx context.Context, id string, dto *entity.StockSessionDto, actorID string) (*entity.StockSessionDto, error) {
	existing, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing.Status == entity.StockSessionStatusClosed {
		return nil, status.New(status.BadRequest, errors.New("session already closed"))
	}

	dto.ID = id
	dto.EmployeeID = existing.EmployeeID
	dto.Date = existing.Date
	if dto.Notes == "" {
		dto.Notes = existing.Notes
	}

	if err := s.hydrateItemSnapshots(ctx, dto); err != nil {
		return nil, err
	}

	// Validation: at close, (cashSoldQty + cashlessSoldQty) +
	// returnQty cannot exceed the morning count we persisted at
	// open time. We look the original outQty up from
	// `existing.Items` by itemId, because the close wire payload
	// only carries the two split sold columns + returnQty (not
	// outQty).
	morningOut := make(map[string]int, len(existing.Items))
	morningItemName := make(map[string]string, len(existing.Items))
	for _, m := range existing.Items {
		morningOut[m.ItemID] = m.OutQty
		if m.Item != nil {
			morningItemName[m.ItemID] = m.Item.Name
		} else {
			morningItemName[m.ItemID] = m.ItemID
		}
	}
	log.WithContext(ctx).Infof("morningOut %+v", morningOut)
	for i, it := range dto.Items {
		// Check whether this item was loaded into the morning session.
		// Items that weren't — typically variants added on the close
		// form after the driver already left — have no morning cap,
		// so we can't validate soldQty + returnQty against anything.
		// We still persist them as new stock_session_item rows with
		// OutQty = 0 (the morning count doesn't apply) and the admin-
		// provided cash/cashless/returnQty. repo.Update replaces all
		// children atomically, so the new row is created as part of
		// the close write without needing an explicit INSERT here.
		if _, ok := morningOut[it.ItemID]; !ok {
			itemName := morningItemName[it.ItemID]
			log.WithContext(ctx).Infof(
				"[stock-session/close] inserting new stock_session_item: item=%s (id=%s) cashSoldQty=%d cashlessSoldQty=%d returnQty=%d outQty=0",
				itemName, it.ItemID, it.CashSoldQty, it.CashlessSoldQty, it.ReturnQty,
			)
			// OutQty stays 0 because this item was not part of the
			// morning session — there is no morning count to track.
			// The split Cash/Cashless + ReturnQty are kept exactly
			// as the admin typed.
			dto.Items[i].OutQty = 0
			continue
		}
		originalOut := morningOut[it.ItemID]
		dto.Items[i].OutQty = originalOut
		soldTotal := it.CashSoldQty + it.CashlessSoldQty
		if soldTotal+it.ReturnQty > originalOut {
			itemName := morningItemName[it.ItemID]
			msg := fmt.Sprintf(
				"(cashSoldQty + cashlessSoldQty) + returnQty cannot exceed morning outQty: item=%s (id=%s) cashSoldQty=%d cashlessSoldQty=%d returnQty=%d morningOutQty=%d total=%d",
				itemName, it.ItemID, it.CashSoldQty, it.CashlessSoldQty, it.ReturnQty, originalOut, soldTotal+it.ReturnQty,
			)
			log.WithContext(ctx).Warnf("[stock-session/close] validation failed: %s", msg)
			return nil, status.New(status.BadRequest, errors.New(msg))
		}
		if it.ReturnQty < 0 || it.CashSoldQty < 0 || it.CashlessSoldQty < 0 {
			itemName := morningItemName[it.ItemID]
			msg := fmt.Sprintf(
				"qty cannot be negative: item=%s (id=%s) cashSoldQty=%d cashlessSoldQty=%d returnQty=%d",
				itemName, it.ItemID, it.CashSoldQty, it.CashlessSoldQty, it.ReturnQty,
			)
			log.WithContext(ctx).Warnf("[stock-session/close] validation failed: %s", msg)
			return nil, status.New(status.BadRequest, errors.New(msg))
		}
	}

	dto.OrganizationID = existing.OrganizationID
	dto.Status = entity.StockSessionStatusClosed
	dto.OpenedAt = existing.OpenedAt
	dto.CreatedBy = existing.CreatedBy
	now := time.Now()
	dto.ClosedAt = &now
	dto.RecomputeTotals()
	s.resolveAndApplySalary(ctx, dto)

	result, err := s.repo.Update(ctx, dto)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) FindAll(ctx context.Context, req *entity.StockSessionFindAllRequest) (*pagination.ResultPagination, error) {
	if req.FindAllRequest.OrganizationData.ID == "" {
		req.FindAllRequest.OrganizationData.ID = shared.GetOrganization(ctx).ID
	}
	return s.repo.FindAll(ctx, req)
}

// hydrateItemSnapshots loads only the price/cost/commission columns
// of the items referenced by the incoming payload and writes them
// into each row's SellingPriceSnapshot / CostPriceSnapshot /
// CommissionSnapshot fields. It deliberately does NOT populate the
// nested Item DTO — that would later be picked up by ToModel() and
// cause GORM to auto-insert an item row on the stock_session_item
// write. The contract here is: stock-session open only references
// existing items by id, never creates them.
func (s *service) hydrateItemSnapshots(ctx context.Context, dto *entity.StockSessionDto) error {
	if len(dto.Items) == 0 {
		return nil
	}
	ids := make([]string, 0, len(dto.Items))
	for _, it := range dto.Items {
		if it.ItemID != "" {
			ids = append(ids, it.ItemID)
		}
	}
	if len(ids) == 0 {
		return nil
	}
	// Only pull the columns we actually need; keep payloads tiny and
	// avoid drifting the in-memory model into the dto.
	var rows []struct {
		ID        string
		Price     float64
		CostPrice float64
		Commision float64
	}
	if err := s.db.Model(&model.Item{}).
		Select("id", "price", "cost_price", "commision").
		Where("id IN ?", ids).
		Scan(&rows).Error; err != nil {
		return err
	}
	byID := make(map[string]struct {
		Price     float64
		CostPrice float64
		Commision float64
	}, len(rows))
	for _, r := range rows {
		byID[r.ID] = struct {
			Price     float64
			CostPrice float64
			Commision float64
		}{Price: r.Price, CostPrice: r.CostPrice, Commision: r.Commision}
	}
	for i := range dto.Items {
		snap, ok := byID[dto.Items[i].ItemID]
		if !ok {
			return status.New(status.BadRequest, errors.New("item not found: "+dto.Items[i].ItemID))
		}
		// Defensive: even if the request body carries a nested Item
		// (older clients), drop it here so it can never leak through
		// to GORM on save.
		dto.Items[i].Item = nil
		if dto.Items[i].SellingPriceSnapshot == 0 {
			dto.Items[i].SellingPriceSnapshot = snap.Price
		}
		if dto.Items[i].CostPriceSnapshot == 0 {
			dto.Items[i].CostPriceSnapshot = snap.CostPrice
		}
		if dto.Items[i].CommissionSnapshot == 0 {
			dto.Items[i].CommissionSnapshot = snap.Commision
		}
	}
	return nil
}

// resolveAndApplySalary looks up the salary_component rows that
// apply to the driver's company, then asks the DTO to compute the
// per-component amounts based on the session's TotalItems.
//
// Lookup path:
//
//	employee_id → admin_company.admin_id → company.id
//	                                         → salary_component.company_id
//
// The salary breakdown is recomputed on every write (Open / Update /
// Close) so reports don't need to re-derive it, and so changes to
// the master salary_component table flow through to past sessions
// the next time they are touched.
//
// Errors are logged but non-fatal — a missing company or empty
// salary bands just produce a 0-amount breakdown for that session.
func (s *service) resolveAndApplySalary(ctx context.Context, dto *entity.StockSessionDto) {
	if dto.EmployeeID == "" {
		return
	}
	var companyID string
	err := s.db.Model(&model.AdminCompany{}).
		Select("company_id").
		Where("admin_id = ?", dto.EmployeeID).
		Scan(&companyID).Error
	if err != nil {
		log.WithContext(ctx).Warnf(
			"[stock-session] salary lookup failed (admin_company): employee=%s err=%v",
			dto.EmployeeID, err,
		)
		return
	}
	if companyID == "" {
		// No company binding yet — leave the salary fields at 0.
		log.WithContext(ctx).Warnf(
			"[stock-session] salary lookup failed (admin_company): employee=%s has no company binding",
			dto.EmployeeID,
		)
		return
	}
	var rows []model.SalaryComponent
	if err := s.db.Where("company_id = ?", companyID).Find(&rows).Error; err != nil {
		log.WithContext(ctx).Warnf(
			"[stock-session] salary lookup failed (salary_component): company=%s err=%v",
			companyID, err,
		)
		return
	}
	components := make([]entity.SalaryComponentDto, 0, len(rows))
	for i := range rows {
		components = append(components, *entity.NewSalaryComponentDtoFromModel(&rows[i]))
	}
	dto.RecomputeSalary(components, dto.TotalItems)
}

// func (s *service) appendLog(ctx context.Context, sessionID, adminID, action, detail string) {
// 	logEntry := &model.SessionLog{
// 		SessionID: sessionID,
// 		Action:    action,
// 		AdminID:   adminID,
// 		Detail:    detail,
// 	}
// 	if err := s.db.Create(logEntry).Error; err != nil {
// 		log.WithContext(ctx).Errorf("failed to write session log: %v", err)
// 	}
// }
