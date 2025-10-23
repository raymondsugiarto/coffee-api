package benefit_type

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, dto *entity.BenefitTypeDto) (*entity.BenefitTypeDto, error)
	FindByID(ctx context.Context, id string) (*entity.BenefitTypeDto, error)
	FindAll(ctx context.Context, req *entity.BenefitTypeFindAllRequest) (*pagination.ResultPagination, error)
	Update(ctx context.Context, dto *entity.BenefitTypeDto) (*entity.BenefitTypeDto, error)
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, dto *entity.BenefitTypeDto) (*entity.BenefitTypeDto, error) {
	m := dto.ToModel()
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	return new(entity.BenefitTypeDto).FromModel(m), nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*entity.BenefitTypeDto, error) {
	var m model.BenefitType
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return new(entity.BenefitTypeDto).FromModel(&m), nil
}

func (r *repository) FindAll(ctx context.Context, req *entity.BenefitTypeFindAllRequest) (*pagination.ResultPagination, error) {
	var m []model.BenefitType = make([]model.BenefitType, 0)

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		return r.db.Model(&model.BenefitType{})
	}, &pagination.TableRequest{
		Request:       req,
		QueryField:    []string{},
		Data:          &m,
		AllowedFields: []string{"name", "status", "created_at"},
		MapFields: map[string]string{
			"name":       "name",
			"status":     "status",
			"created_at": "created_at",
		},
	})
	if err != nil {
		return nil, err
	}

	result := dataTable.(*pagination.ResultPagination)
	results := result.Data.(*[]model.BenefitType)
	var data []*entity.BenefitTypeDto = make([]*entity.BenefitTypeDto, 0)
	for _, m := range *results {
		data = append(data, new(entity.BenefitTypeDto).FromModel(&m))
	}
	return &pagination.ResultPagination{
		Data:        data,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}

func (r *repository) Update(ctx context.Context, dto *entity.BenefitTypeDto) (*entity.BenefitTypeDto, error) {
	m := dto.ToModel()
	if err := r.db.WithContext(ctx).Where("id = ?", dto.ID).Updates(m).Error; err != nil {
		return nil, err
	}
	return new(entity.BenefitTypeDto).FromModel(m), nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.BenefitType{}).Error
}
