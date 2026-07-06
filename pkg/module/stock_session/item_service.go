package stocksession

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

// ItemService provides read access to existing `item` and `item_category` tables
// for the stock-session flow (driver picker).
type ItemService interface {
	FindAllItems(ctx context.Context, req *entity.ItemFindAllRequest) (*pagination.ResultPagination, error)
	GetItem(ctx context.Context, id string) (*entity.ItemDto, error)
}

type itemService struct {
	db *gorm.DB
}

func NewItemService(db *gorm.DB) ItemService {
	return &itemService{db: db}
}

func (s *itemService) FindAllItems(ctx context.Context, req *entity.ItemFindAllRequest) (*pagination.ResultPagination, error) {
	var m []model.Item = make([]model.Item, 0)
	tbl := pagination.NewTable(s.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		q := s.db.Model(&model.Item{})
		if req.IsActive != nil {
			q = q.Where("is_active = ?", *req.IsActive)
		} else {
			q = q.Where("is_active = ?", true)
		}
		if req.CategoryID != "" {
			q = q.Where("category_id = ?", req.CategoryID)
		}
		if req.Query != "" {
			like := "%" + req.Query + "%"
			q = q.Where("name ILIKE ? OR sku ILIKE ? OR code ILIKE ?", like, like, like)
		}
		return q
	}, &pagination.TableRequest{
		Request:       req,
		QueryField:    []string{"name", "sku", "code"},
		Data:          &m,
		AllowedFields: []string{"name", "price", "created_at"},
	})
	if err != nil {
		return nil, err
	}
	result := dataTable.(*pagination.ResultPagination)
	results := result.Data.(*[]model.Item)
	data := make([]*entity.ItemDto, 0, len(*results))
	for i := range *results {
		data = append(data, entity.NewItemDtoFromModel(&(*results)[i]))
	}
	return &pagination.ResultPagination{
		Data:        data,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}

func (s *itemService) GetItem(ctx context.Context, id string) (*entity.ItemDto, error) {
	var m *model.Item
	if err := s.db.Where("id = ?", id).Preload("Category").First(&m).Error; err != nil {
		return nil, err
	}
	return entity.NewItemDtoFromModel(m), nil
}
