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
	m := dto.ToModel()
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return entity.NewStockSessionDtoFromModel(m), nil
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
