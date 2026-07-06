package stocksession

import (
	"context"
	"errors"
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

	result, err := s.repo.Create(ctx, dto)
	if err != nil {
		return nil, err
	}

	s.appendLog(ctx, result.ID, actorID, entity.SessionActionOpen, "Session opened")
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

	result, err := s.repo.Update(ctx, dto)
	if err != nil {
		return nil, err
	}
	s.appendLog(ctx, dto.ID, actorID, entity.SessionActionUpdate, "Session updated")
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

	// Validation: return_qty cannot exceed out_qty
	for _, it := range dto.Items {
		if it.ReturnQty > it.OutQty {
			return nil, status.New(status.BadRequest, errors.New("returnQty cannot exceed outQty"))
		}
	}

	dto.OrganizationID = existing.OrganizationID
	dto.Status = entity.StockSessionStatusClosed
	dto.OpenedAt = existing.OpenedAt
	dto.CreatedBy = existing.CreatedBy
	now := time.Now()
	dto.ClosedAt = &now
	dto.RecomputeTotals()

	result, err := s.repo.Update(ctx, dto)
	if err != nil {
		return nil, err
	}
	s.appendLog(ctx, id, actorID, entity.SessionActionClose, "Session closed")
	return result, nil
}

func (s *service) FindAll(ctx context.Context, req *entity.StockSessionFindAllRequest) (*pagination.ResultPagination, error) {
	if req.FindAllRequest.OrganizationData.ID == "" {
		req.FindAllRequest.OrganizationData.ID = shared.GetOrganization(ctx).ID
	}
	return s.repo.FindAll(ctx, req)
}

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
	var items []model.Item
	if err := s.db.Where("id IN ?", ids).Find(&items).Error; err != nil {
		return err
	}
	byID := make(map[string]*model.Item, len(items))
	for i := range items {
		byID[items[i].ID] = &items[i]
	}
	for i := range dto.Items {
		it, ok := byID[dto.Items[i].ItemID]
		if !ok {
			return status.New(status.BadRequest, errors.New("item not found: "+dto.Items[i].ItemID))
		}
		dto.Items[i].Item = entity.NewItemDtoFromModel(it)
		if dto.Items[i].SellingPriceSnapshot == 0 {
			dto.Items[i].SellingPriceSnapshot = it.Price
		}
		if dto.Items[i].CostPriceSnapshot == 0 {
			dto.Items[i].CostPriceSnapshot = it.CostPrice
		}
	}
	return nil
}

func (s *service) appendLog(ctx context.Context, sessionID, adminID, action, detail string) {
	logEntry := &model.SessionLog{
		SessionID: sessionID,
		Action:    action,
		AdminID:   adminID,
		Detail:    detail,
	}
	if err := s.db.Create(logEntry).Error; err != nil {
		log.WithContext(ctx).Errorf("failed to write session log: %v", err)
	}
}
