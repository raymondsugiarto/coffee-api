package admin

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, dto *entity.AdminDto, tx *gorm.DB) (*entity.AdminDto, error)
	Update(ctx context.Context, dto *entity.AdminDto) (*entity.AdminDto, error)
	FindAll(ctx context.Context, req *entity.FindAllRequest) (*pagination.ResultPagination, error)
	FindByUserID(ctx context.Context, id string) (*entity.AdminDto, error)
	CreateAdminCompany(ctx context.Context, dto *entity.CreateAdminCompany, cb func(tx *gorm.DB) error) (*entity.CreateAdminCompany, error)
	FindAllByCompanyID(ctx context.Context, companyID string, req *entity.FindAllRequest) (*pagination.ResultPagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, dto *entity.AdminDto, tx *gorm.DB) (*entity.AdminDto, error) {
	db := r.db
	if tx != nil {
		db = tx
	}
	m := dto.ToModel()
	err := db.Create(m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.AdminDto).FromModel(m), nil
}

func (r *repository) Update(ctx context.Context, dto *entity.AdminDto) (*entity.AdminDto, error) {
	m := dto.ToModel()
	err := r.db.Model(m).Where("user_id = ?", m.UserID).Updates(m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.AdminDto).FromModel(m), nil
}

func (r *repository) FindAll(ctx context.Context, req *entity.FindAllRequest) (*pagination.ResultPagination, error) {
	var m []model.Admin = make([]model.Admin, 0)

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		return r.db.Model(&model.Admin{}).
			Where("deleted_at is null")
	}, &pagination.TableRequest{
		Request:       req,
		QueryField:    []string{},
		Data:          &m,
		AllowedFields: []string{""},
	})
	if err != nil {
		return nil, err
	}

	result := dataTable.(*pagination.ResultPagination)
	results := result.Data.(*[]model.Admin)
	var data []*entity.AdminDto = make([]*entity.AdminDto, 0)
	for _, m := range *results {
		data = append(data, new(entity.AdminDto).FromModel(&m))
	}
	return &pagination.ResultPagination{
		Data:        data,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}

func (r *repository) FindByUserID(ctx context.Context, id string) (*entity.AdminDto, error) {
	var m *model.Admin
	err := r.db.Where("user_id = ?", id).First(&m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.AdminDto).FromModel(m), nil
}

func (r *repository) CreateAdminCompany(ctx context.Context, dto *entity.CreateAdminCompany, cb func(tx *gorm.DB) error) (*entity.CreateAdminCompany, error) {
	m := dto.ToModel()
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(m).Error; err != nil {
			return err
		}
		dto.ID = m.ID
		dto.UserID = m.UserID
		if cb != nil {
			if err := cb(tx); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return new(entity.CreateAdminCompany).FromModel(m), nil
}

func (r *repository) FindAllByCompanyID(ctx context.Context, companyID string, req *entity.FindAllRequest) (*pagination.ResultPagination, error) {
	var m []model.Admin = make([]model.Admin, 0)

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		return r.db.Model(&model.Admin{}).
			Where("company_id = ?", companyID)
	}, &pagination.TableRequest{
		Request:       req,
		QueryField:    []string{"first_name", "last_name"},
		Data:          &m,
		AllowedFields: []string{""},
	})
	if err != nil {
		return nil, err
	}

	result := dataTable.(*pagination.ResultPagination)
	results := result.Data.(*[]model.Admin)

	var data []*entity.AdminDto = make([]*entity.AdminDto, 0)
	for _, m := range *results {
		data = append(data, new(entity.AdminDto).FromModel(&m))
	}

	return &pagination.ResultPagination{
		Data:        data,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}
