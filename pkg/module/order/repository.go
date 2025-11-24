package order

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, role *entity.OrderDto) (*entity.OrderDto, error)
	Get(ctx context.Context, id string) (*entity.OrderDto, error)
	Update(ctx context.Context, role *entity.OrderDto) (*entity.OrderDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.OrderFindAllRequest) (*pagination.ResultPagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, role *entity.OrderDto) (*entity.OrderDto, error) {
	m := role.ToModel()
	err := r.db.Create(m).Error
	if err != nil {
		return nil, err
	}
	return entity.NewOrderDtoFromModel(m), nil
}

func (r *repository) Get(ctx context.Context, id string) (*entity.OrderDto, error) {
	var m *model.Order
	err := r.db.Where("id = ?", id).Preload("RoleParent").First(&m).Error
	if err != nil {
		return nil, err
	}
	return entity.NewOrderDtoFromModel(m), nil
}

func (r *repository) Update(ctx context.Context, role *entity.OrderDto) (*entity.OrderDto, error) {
	err := r.db.Save(role.ToModel()).Error
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	err := r.db.Where("id = ?", id).Delete(&model.Order{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) FindAll(ctx context.Context, req *entity.OrderFindAllRequest) (*pagination.ResultPagination, error) {
	var m []model.Order = make([]model.Order, 0)

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		return r.db.Model(&model.Order{}).Where("admin_id = ? AND company_id = ?", req.AdminID, req.CompanyID)
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
	results := result.Data.(*[]model.Order)
	var data []*entity.OrderDto = make([]*entity.OrderDto, 0)
	for _, m := range *results {
		data = append(data, entity.NewOrderDtoFromModel(&m))
	}
	return &pagination.ResultPagination{
		Data:        data,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}
