package redeem

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, dto *entity.RedeemDto) (*entity.RedeemDto, error)
	CreateWithTx(ctx context.Context, tx *gorm.DB, dto *entity.RedeemDto) (*entity.RedeemDto, error)
	Get(ctx context.Context, id string) (*entity.RedeemDto, error)
	Update(ctx context.Context, dto *entity.RedeemDto) (*entity.RedeemDto, error)
	UpdateWithTx(ctx context.Context, tx *gorm.DB, dto *entity.RedeemDto) (*entity.RedeemDto, error)
	Delete(ctx context.Context, id string) error
	DeleteWithTx(ctx context.Context, tx *gorm.DB, id string) error
	FindAll(ctx context.Context, req *entity.RedeemFindAllRequest) (*pagination.ResultPagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, dto *entity.RedeemDto) (*entity.RedeemDto, error) {
	m := dto.ToModel()
	err := r.db.Create(m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.RedeemDto).FromModel(m), nil
}

func (r *repository) CreateWithTx(ctx context.Context, tx *gorm.DB, dto *entity.RedeemDto) (*entity.RedeemDto, error) {
	m := dto.ToModel()
	err := tx.Create(m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.RedeemDto).FromModel(m), nil
}

func (r *repository) Get(ctx context.Context, id string) (*entity.RedeemDto, error) {
	var m *model.Redeem
	err := r.db.Preload("Customer").
		Preload("Reward").
		Preload("Province").
		Preload("Regency").
		Preload("District").
		Preload("Village").
		Where("id = ?", id).First(&m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.RedeemDto).FromModel(m), nil
}

func (r *repository) Update(ctx context.Context, dto *entity.RedeemDto) (*entity.RedeemDto, error) {
	err := r.db.Updates(dto.ToModel()).Error
	if err != nil {
		return nil, err
	}
	return dto, nil
}

func (r *repository) UpdateWithTx(ctx context.Context, tx *gorm.DB, dto *entity.RedeemDto) (*entity.RedeemDto, error) {
	err := tx.Updates(dto.ToModel()).Error
	if err != nil {
		return nil, err
	}
	return dto, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	err := r.db.Where("id = ?", id).Delete(&model.Redeem{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) DeleteWithTx(ctx context.Context, tx *gorm.DB, id string) error {
	err := tx.Where("id = ?", id).Delete(&model.Redeem{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) FindAll(ctx context.Context, req *entity.RedeemFindAllRequest) (*pagination.ResultPagination, error) {
	var m []model.Redeem = make([]model.Redeem, 0)

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i any) *gorm.DB {
		query := r.db.Model(&model.Redeem{}).
			Preload("Customer").
			Preload("Reward").
			Preload("Province").
			Preload("Regency").
			Preload("District").
			Preload("Village").
			Where("deleted_at is null")

		// Apply filters
		if req.CustomerID != "" {
			query = query.Where("customer_id = ?", req.CustomerID)
		}
		if req.Status != "" {
			query = query.Where("status = ?", req.Status)
		}

		return query.Order("created_at desc")
	}, &pagination.TableRequest{
		Request:       &req.FindAllRequest,
		QueryField:    []string{"redemption_code"},
		Data:          &m,
		AllowedFields: []string{"redemption_code"},
	})
	if err != nil {
		return nil, err
	}

	result := dataTable.(*pagination.ResultPagination)
	results := result.Data.(*[]model.Redeem)
	var data []*entity.RedeemDto = make([]*entity.RedeemDto, 0)
	for _, m := range *results {
		data = append(data, new(entity.RedeemDto).FromModel(&m))
	}
	return &pagination.ResultPagination{
		Data:        data,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}
