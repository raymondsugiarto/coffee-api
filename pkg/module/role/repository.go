package role

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, role *entity.RoleDto) (*entity.RoleDto, error)
	Get(ctx context.Context, id string) (*entity.RoleDto, error)
	Update(ctx context.Context, role *entity.RoleDto) (*entity.RoleDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.FindAllRequest) (*pagination.ResultPagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, role *entity.RoleDto) (*entity.RoleDto, error) {
	m := role.ToModel()
	err := r.db.Create(m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.RoleDto).FromModel(m), nil
}

func (r *repository) Get(ctx context.Context, id string) (*entity.RoleDto, error) {
	var m *model.Role
	err := r.db.Where("id = ?", id).Preload("RoleParent").First(&m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.RoleDto).FromModel(m), nil
}

func (r *repository) Update(ctx context.Context, role *entity.RoleDto) (*entity.RoleDto, error) {
	err := r.db.Save(role.ToModel()).Error
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	err := r.db.Where("id = ?", id).Delete(&model.Role{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) FindAll(ctx context.Context, req *entity.FindAllRequest) (*pagination.ResultPagination, error) {
	var m []model.Role = make([]model.Role, 0)

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		return r.db.Model(&model.Role{}).
			Preload("RoleParent").
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
	results := result.Data.(*[]model.Role)
	var data []*entity.RoleDto = make([]*entity.RoleDto, 0)
	for _, m := range *results {
		data = append(data, new(entity.RoleDto).FromModel(&m))
	}
	return &pagination.ResultPagination{
		Data:        data,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}
