package investment

import (
	"context"

	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/investment"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type CallbackCreate func(ctx context.Context, db *gorm.DB, dto *entity.InvestmentDto) error

type Repository interface {
	Create(ctx context.Context, dto *entity.InvestmentDto, cb CallbackCreate) (*entity.InvestmentDto, error)
	Get(ctx context.Context, id string) (*entity.InvestmentDto, error)
	FindByCode(ctx context.Context, code string) (*entity.InvestmentDto, error)
	Update(ctx context.Context, dto *entity.InvestmentDto) (*entity.InvestmentDto, error)
	UpdateWithTx(ctx context.Context, tx *gorm.DB, dto *entity.InvestmentDto) (*entity.InvestmentDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.InvestmentFindAllRequest) (*pagination.ResultPagination, error)
	UpdateStatusWaitingVerification(ctx context.Context, dto *entity.InvestmentDto) (*entity.InvestmentDto, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, dto *entity.InvestmentDto, cb CallbackCreate) (*entity.InvestmentDto, error) {
	m := dto.ToModel()
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(m).Error; err != nil {
			return err
		}
		dto.ID = m.ID
		if err := cb(ctx, tx, dto); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return new(entity.InvestmentDto).FromModel(m), nil
}

func (r *repository) Get(ctx context.Context, id string) (*entity.InvestmentDto, error) {
	var m *model.Investment
	err := r.db.Where("id = ?", id).
		Preload("InvestmentPayments.Bank").
		Preload("InvestmentPayments").
		Preload("InvestmentItems").
		Preload("Customer").
		First(&m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.InvestmentDto).FromModel(m), nil
}

func (r *repository) FindByCode(ctx context.Context, code string) (*entity.InvestmentDto, error) {
	var m *model.Investment
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.InvestmentDto).FromModel(m), nil
}

func (r *repository) Update(ctx context.Context, dto *entity.InvestmentDto) (*entity.InvestmentDto, error) {
	err := r.db.Updates(dto.ToModel()).Where("id = ? ", dto.ID).Error
	if err != nil {
		return nil, err
	}
	return dto, nil
}

func (r *repository) UpdateWithTx(ctx context.Context, tx *gorm.DB, dto *entity.InvestmentDto) (*entity.InvestmentDto, error) {
	err := tx.Model(&model.Investment{}).Where("id = ?", dto.ID).Updates(dto.ToModel()).Error
	if err != nil {
		return nil, err
	}
	return dto, nil
}

func (r *repository) UpdateStatusWaitingVerification(ctx context.Context, dto *entity.InvestmentDto) (*entity.InvestmentDto, error) {
	err := r.db.Model(&model.Investment{}).Where("id = ?", dto.ID).Updates(map[string]interface{}{
		"status": model.InvestmentStatusRequest,
	}).Error
	if err != nil {
		return nil, err
	}
	return dto, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	err := r.db.Where("id = ?", id).Delete(&model.Investment{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) FindAll(ctx context.Context, req *entity.InvestmentFindAllRequest) (*pagination.ResultPagination, error) {
	var m []model.Investment = make([]model.Investment, 0)

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		q := r.db.Model(&model.Investment{}).
			Preload("Company").
			Preload("Customer")

		if req.IncludePayments {
			q = q.Preload("InvestmentPayments.Bank")
			q = q.Preload("InvestmentPayments")
		}

		q = q.Where("status != ?", model.InvestmentStatusCreated).Order("created_at DESC")
		return q
	}, &pagination.TableRequest{
		Request:       req,
		QueryField:    []string{},
		Data:          &m,
		AllowedFields: []string{"customer_id", "company_id", "type", "created_at"},
	})
	if err != nil {
		return nil, err
	}

	result := dataTable.(*pagination.ResultPagination)
	results := result.Data.(*[]model.Investment)
	var data []*entity.InvestmentDto = make([]*entity.InvestmentDto, 0)
	for _, m := range *results {
		data = append(data, new(entity.InvestmentDto).FromModel(&m))
	}
	return &pagination.ResultPagination{
		Data:        data,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}
