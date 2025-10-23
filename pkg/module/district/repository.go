package district

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	FindAll(ctx context.Context, req *entity.DistrictFindAllRequest) (*pagination.ResultPagination, error)
	FindByID(ctx context.Context, id string) (*entity.DistrictDto, error)
	FindByCode(ctx context.Context, code string) (*entity.DistrictDto, error)
	FindByRegencyID(ctx context.Context, regencyID string) ([]*entity.DistrictDto, error)
	FindByRegencyCode(ctx context.Context, regencyCode string) ([]*entity.DistrictDto, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindAll(ctx context.Context, req *entity.DistrictFindAllRequest) (*pagination.ResultPagination, error) {
	var m []model.District

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		q := r.db.Model(&model.District{}).
			Preload("Regency.Province")

		if req.RegencyCode != "" {
			q = q.Joins("JOIN regency ON regency.id = district.regency_id").
				Where("regency.code = ?", req.RegencyCode)
		}

		return q
	}, &pagination.TableRequest{
		Request:       req,
		QueryField:    []string{"name", "code"},
		Data:          &m,
		AllowedFields: []string{"code", "regency_id"},
	})
	if err != nil {
		return nil, err
	}

	result := dataTable.(*pagination.ResultPagination)
	results := result.Data.(*[]model.District)
	var data []*entity.DistrictDto = make([]*entity.DistrictDto, 0)
	for _, m := range *results {
		data = append(data, new(entity.DistrictDto).FromModel(&m))
	}
	return &pagination.ResultPagination{
		Data:        data,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*entity.DistrictDto, error) {
	var m model.District
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Preload("Regency.Province").
		First(&m).Error
	if err != nil {
		return nil, err
	}

	return new(entity.DistrictDto).FromModel(&m), nil
}

func (r *repository) FindByCode(ctx context.Context, code string) (*entity.DistrictDto, error) {
	var m model.District
	err := r.db.WithContext(ctx).
		Where("code = ?", code).
		Preload("Regency.Province").
		First(&m).Error
	if err != nil {
		return nil, err
	}

	return new(entity.DistrictDto).FromModel(&m), nil
}

func (r *repository) FindByRegencyID(ctx context.Context, regencyID string) ([]*entity.DistrictDto, error) {
	var m []model.District
	err := r.db.WithContext(ctx).
		Joins("JOIN regency ON regency.id = district.regency_id").
		Where("regency.id = ?", regencyID).
		Preload("Regency.Province").
		Find(&m).Error
	if err != nil {
		return nil, err
	}

	var d []*entity.DistrictDto
	for _, district := range m {
		d = append(d, new(entity.DistrictDto).FromModel(&district))
	}
	return d, nil
}

func (r *repository) FindByRegencyCode(ctx context.Context, regencyCode string) ([]*entity.DistrictDto, error) {
	var m []model.District
	err := r.db.WithContext(ctx).
		Joins("JOIN regency ON regency.id = district.regency_id").
		Where("regency.code = ?", regencyCode).
		Preload("Regency.Province").
		Find(&m).Error
	if err != nil {
		return nil, err
	}

	var d []*entity.DistrictDto
	for _, district := range m {
		d = append(d, new(entity.DistrictDto).FromModel(&district))
	}
	return d, nil
}
