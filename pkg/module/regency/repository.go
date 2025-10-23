package regency

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	FindAll(ctx context.Context, req *entity.RegencyFindAllRequest) (*pagination.ResultPagination, error)
	FindByID(ctx context.Context, id string) (*entity.RegencyDto, error)
	FindByCode(ctx context.Context, code string) (*entity.RegencyDto, error)
	FindByProvinceID(ctx context.Context, provinceID string) ([]*entity.RegencyDto, error)
	FindByProvinceCode(ctx context.Context, provinceCode string) ([]*entity.RegencyDto, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindAll(ctx context.Context, req *entity.RegencyFindAllRequest) (*pagination.ResultPagination, error) {
	var m []model.Regency

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		q := r.db.Model(&model.Regency{}).
			Preload("Province")

		if req.ProvinceCode != "" {
			q = q.Joins("JOIN province ON province.id = regency.province_id").
				Where("province.code = ?", req.ProvinceCode)
		}

		return q
	}, &pagination.TableRequest{
		Request:       req,
		QueryField:    []string{"name", "code"},
		Data:          &m,
		AllowedFields: []string{"code", "province_id"},
	})
	if err != nil {
		return nil, err
	}

	result := dataTable.(*pagination.ResultPagination)
	results := result.Data.(*[]model.Regency)
	var data []*entity.RegencyDto = make([]*entity.RegencyDto, 0)
	for _, m := range *results {
		data = append(data, new(entity.RegencyDto).FromModel(&m))
	}
	return &pagination.ResultPagination{
		Data:        data,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*entity.RegencyDto, error) {
	var m model.Regency
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Preload("Province").
		First(&m).Error
	if err != nil {
		return nil, err
	}

	return new(entity.RegencyDto).FromModel(&m), nil
}

func (r *repository) FindByCode(ctx context.Context, code string) (*entity.RegencyDto, error) {
	var m model.Regency
	err := r.db.WithContext(ctx).
		Where("code = ?", code).
		Preload("Province").
		First(&m).Error
	if err != nil {
		return nil, err
	}

	return new(entity.RegencyDto).FromModel(&m), nil
}

func (r *repository) FindByProvinceID(ctx context.Context, provinceID string) ([]*entity.RegencyDto, error) {
	var m []model.Regency
	err := r.db.WithContext(ctx).
		Joins("JOIN province ON province.id = regency.province_id").
		Where("province.id = ?", provinceID).
		Preload("Province").
		Find(&m).Error
	if err != nil {
		return nil, err
	}

	var d []*entity.RegencyDto
	for _, regency := range m {
		d = append(d, new(entity.RegencyDto).FromModel(&regency))
	}
	return d, nil
}

func (r *repository) FindByProvinceCode(ctx context.Context, provinceCode string) ([]*entity.RegencyDto, error) {
	var m []model.Regency
	err := r.db.WithContext(ctx).
		Joins("JOIN province ON province.id = regency.province_id").
		Where("province.code = ?", provinceCode).
		Preload("Province").
		Find(&m).Error
	if err != nil {
		return nil, err
	}

	var d []*entity.RegencyDto
	for _, regency := range m {
		d = append(d, new(entity.RegencyDto).FromModel(&regency))
	}
	return d, nil
}
