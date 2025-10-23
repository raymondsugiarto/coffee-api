package village

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	FindAll(ctx context.Context, req *entity.VillageFindAllRequest) (*pagination.ResultPagination, error)
	FindByID(ctx context.Context, id string) (*entity.VillageDto, error)
	FindByCode(ctx context.Context, code string) (*entity.VillageDto, error)
	FindByDistrictID(ctx context.Context, districtId string) ([]*entity.VillageDto, error)
	FindByDistrictCode(ctx context.Context, districtCode string) ([]*entity.VillageDto, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindAll(ctx context.Context, req *entity.VillageFindAllRequest) (*pagination.ResultPagination, error) {
	var m []model.Village

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		q := r.db.Model(&model.Village{}).
			Preload("District.Regency.Province")

		if req.DistrictCode != "" {
			q = q.Joins("JOIN district ON district.id = village.district_id").
				Where("district.code = ?", req.DistrictCode)
		}

		return q
	}, &pagination.TableRequest{
		Request:       req,
		QueryField:    []string{"name", "code"},
		Data:          &m,
		AllowedFields: []string{"code", "district_id"},
	})
	if err != nil {
		return nil, err
	}

	result := dataTable.(*pagination.ResultPagination)
	results := result.Data.(*[]model.Village)
	var data []*entity.VillageDto = make([]*entity.VillageDto, 0)
	for _, m := range *results {
		data = append(data, new(entity.VillageDto).FromModel(&m))
	}
	return &pagination.ResultPagination{
		Data:        data,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*entity.VillageDto, error) {
	var m model.Village
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Preload("District.Regency.Province").
		First(&m).Error
	if err != nil {
		return nil, err
	}

	return new(entity.VillageDto).FromModel(&m), nil
}

func (r *repository) FindByCode(ctx context.Context, code string) (*entity.VillageDto, error) {
	var m model.Village
	err := r.db.WithContext(ctx).
		Where("code = ?", code).
		Preload("District.Regency.Province").
		First(&m).Error
	if err != nil {
		return nil, err
	}

	return new(entity.VillageDto).FromModel(&m), nil
}

func (r *repository) FindByDistrictID(ctx context.Context, districtId string) ([]*entity.VillageDto, error) {
	var m []model.Village
	err := r.db.WithContext(ctx).
		Joins("JOIN district ON district.id = village.district_id").
		Where("district.id = ?", districtId).
		Preload("District.Regency.Province").
		Find(&m).Error
	if err != nil {
		return nil, err
	}

	var d []*entity.VillageDto
	for _, village := range m {
		d = append(d, new(entity.VillageDto).FromModel(&village))
	}
	return d, nil
}

func (r *repository) FindByDistrictCode(ctx context.Context, districtCode string) ([]*entity.VillageDto, error) {
	var m []model.Village
	err := r.db.WithContext(ctx).
		Joins("JOIN district ON district.id = village.district_id").
		Where("district.code = ?", districtCode).
		Preload("District.Regency.Province").
		Find(&m).Error
	if err != nil {
		return nil, err
	}

	var d []*entity.VillageDto
	for _, village := range m {
		d = append(d, new(entity.VillageDto).FromModel(&village))
	}
	return d, nil
}
