package country

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	FindAll(ctx context.Context, req *entity.CountryFindAllRequest) (*pagination.ResultPagination, error)
	FindByID(ctx context.Context, id string) (*entity.CountryDto, error)
	FindByCCA2(ctx context.Context, cca2 string) (*entity.CountryDto, error)
	FindByCCA3(ctx context.Context, cca3 string) (*entity.CountryDto, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindAll(ctx context.Context, req *entity.CountryFindAllRequest) (*pagination.ResultPagination, error) {
	var m []model.Country

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		q := r.db.Model(&model.Country{}).
			Where("country.deleted_at is null")
		return q
	}, &pagination.TableRequest{
		Request:       req,
		QueryField:    []string{"name", "cca2", "cca3"},
		Data:          &m,
		AllowedFields: []string{"cca2", "cca3"},
	})
	if err != nil {
		return nil, err
	}

	result := dataTable.(*pagination.ResultPagination)
	results := result.Data.(*[]model.Country)
	var data []*entity.CountryDto = make([]*entity.CountryDto, 0)
	for _, m := range *results {
		data = append(data, new(entity.CountryDto).FromModel(&m))
	}
	return &pagination.ResultPagination{
		Data:        data,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*entity.CountryDto, error) {
	var m model.Country
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&m).Error
	if err != nil {
		return nil, err
	}

	return new(entity.CountryDto).FromModel(&m), nil
}

func (r *repository) FindByCCA2(ctx context.Context, cca2 string) (*entity.CountryDto, error) {
	var m model.Country
	err := r.db.WithContext(ctx).
		Where("cca2 = ?", cca2).
		First(&m).Error
	if err != nil {
		return nil, err
	}

	return new(entity.CountryDto).FromModel(&m), nil
}

func (r *repository) FindByCCA3(ctx context.Context, cca3 string) (*entity.CountryDto, error) {
	var m model.Country
	err := r.db.WithContext(ctx).
		Where("cca3 = ?", cca3).
		First(&m).Error
	if err != nil {
		return nil, err
	}

	return new(entity.CountryDto).FromModel(&m), nil
}
