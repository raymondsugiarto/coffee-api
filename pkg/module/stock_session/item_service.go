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
	GetItemChildren(ctx context.Context, parentIDs []string, includeInactive bool) ([]*entity.ItemDto, error)
	SetItemParent(ctx context.Context, parentID string, childIDs []string) error
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
		// Optional: only show items that have no parent (i.e. parents only).
		// When ParentIDs is empty and parentOnly is desired, callers can
		// rely on GET /api/products/parents instead.
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
	if err := s.db.Where("id = ?", id).Preload("Category").Preload("Parent").First(&m).Error; err != nil {
		return nil, err
	}
	return entity.NewItemDtoFromModel(m), nil
}

// GetItemChildren returns all items whose parent_id matches any of the given
// parentIDs. Used by the close-session UI to expand "Matcha" -> ["Matcha Mango",
// "Matcha Strawberry"] etc. Pass `includeInactive=true` to surface variants the
// admin disabled (rare, but useful for backfills).
func (s *itemService) GetItemChildren(ctx context.Context, parentIDs []string, includeInactive bool) ([]*entity.ItemDto, error) {
	if len(parentIDs) == 0 {
		return []*entity.ItemDto{}, nil
	}
	q := s.db.Model(&model.Item{}).Where("parent_id IN ?", parentIDs)
	if !includeInactive {
		q = q.Where("is_active = ?", true)
	}
	var items []model.Item
	if err := q.Order("name asc").Find(&items).Error; err != nil {
		return nil, err
	}
	out := make([]*entity.ItemDto, 0, len(items))
	for i := range items {
		out = append(out, entity.NewItemDtoFromModel(&items[i]))
	}
	return out, nil
}

// SetItemParent returns nil if every (childID -> parentID) update succeeded.
// Errors out if any of the supplied children / parents don't exist. The
// relation is keyed on UUID (`id`), not on `code` — multiple items can share
// a code (one from seed, others created via UI), so we always look up by id.
func (s *itemService) SetItemParent(ctx context.Context, parentID string, childIDs []string) error {
	if len(childIDs) == 0 {
		return nil
	}
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// Verify parent exists first — clearer failure when caller passes a typo.
	var parent model.Item
	if err := tx.Where("id = ?", parentID).First(&parent).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&model.Item{}).
		Where("id IN ?", childIDs).
		Update("parent_id", parentID).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
