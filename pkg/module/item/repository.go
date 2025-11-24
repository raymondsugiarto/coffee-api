package item

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, role *entity.ItemDto) (*entity.ItemDto, error)
	Get(ctx context.Context, id string) (*entity.ItemDto, error)
	Update(ctx context.Context, role *entity.ItemDto) (*entity.ItemDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.ItemFindAllRequest) (*pagination.ResultPagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, role *entity.ItemDto) (*entity.ItemDto, error) {
	m := role.ToModel()
	err := r.db.Create(m).Error
	if err != nil {
		return nil, err
	}
	return entity.NewItemDtoFromModel(m), nil
}

func (r *repository) Get(ctx context.Context, id string) (*entity.ItemDto, error) {
	var m *model.Item
	err := r.db.Where("id = ?", id).Preload("RoleParent").First(&m).Error
	if err != nil {
		return nil, err
	}
	return entity.NewItemDtoFromModel(m), nil
}

func (r *repository) Update(ctx context.Context, role *entity.ItemDto) (*entity.ItemDto, error) {
	err := r.db.Save(role.ToModel()).Error
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	err := r.db.Where("id = ?", id).Delete(&model.Item{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) FindAll(ctx context.Context, req *entity.ItemFindAllRequest) (*pagination.ResultPagination, error) {
	var m []model.Item = make([]model.Item, 0)

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		return r.db.Model(&model.Item{}).
			Joins("JOIN item_company ON item.id = item_company.item_id").
			Where("item_company.company_id = ?", req.CompanyID)
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
	results := result.Data.(*[]model.Item)
	var data []*entity.ItemDto = make([]*entity.ItemDto, 0)
	for _, m := range *results {
		data = append(data, entity.NewItemDtoFromModel(&m))
	}
	return &pagination.ResultPagination{
		Data:        data,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}
