package province

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	FindAll(ctx context.Context, req *entity.ProvinceFindAllRequest) (*pagination.ResultPagination, error)
	FindByID(ctx context.Context, id string) (*entity.ProvinceDto, error)
	FindByCode(ctx context.Context, code string) (*entity.ProvinceDto, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindAll(ctx context.Context, req *entity.ProvinceFindAllRequest) (*pagination.ResultPagination, error) {
	var m []model.Province

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		q := r.db.Model(&model.Province{}).
			Where("province.deleted_at is null")
		return q
	}, &pagination.TableRequest{
		Request:       req,
		QueryField:    []string{"name", "code"},
		Data:          &m,
		AllowedFields: []string{"code"},
	})
	if err != nil {
		return nil, err
	}

	result := dataTable.(*pagination.ResultPagination)
	results := result.Data.(*[]model.Province)
	var data []*entity.ProvinceDto = make([]*entity.ProvinceDto, 0)
	for _, m := range *results {
		data = append(data, new(entity.ProvinceDto).FromModel(&m))
	}
	return &pagination.ResultPagination{
		Data:        data,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*entity.ProvinceDto, error) {
	var m model.Province
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&m).Error
	if err != nil {
		return nil, err
	}

	return new(entity.ProvinceDto).FromModel(&m), nil
}

func (r *repository) FindByCode(ctx context.Context, code string) (*entity.ProvinceDto, error) {
	var m model.Province
	err := r.db.WithContext(ctx).
		Where("code = ?", code).
		First(&m).Error
	if err != nil {
		return nil, err
	}

	return new(entity.ProvinceDto).FromModel(&m), nil
}
