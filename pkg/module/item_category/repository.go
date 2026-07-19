package itemcategory

import (
	"context"
	"strings"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, dto *entity.ItemCategoryDto) (*entity.ItemCategoryDto, error)
	Get(ctx context.Context, id string) (*entity.ItemCategoryDto, error)
	Update(ctx context.Context, dto *entity.ItemCategoryDto) (*entity.ItemCategoryDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.ItemCategoryFindAllRequest) (*pagination.ResultPagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, dto *entity.ItemCategoryDto) (*entity.ItemCategoryDto, error) {
	m := dto.ToModel()
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return entity.NewItemCategoryDtoFromModel(m), nil
}

func (r *repository) Get(ctx context.Context, id string) (*entity.ItemCategoryDto, error) {
	var m model.ItemCategory
	if err := r.db.Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return entity.NewItemCategoryDtoFromModel(&m), nil
}

func (r *repository) Update(ctx context.Context, dto *entity.ItemCategoryDto) (*entity.ItemCategoryDto, error) {
	if err := r.db.Save(dto.ToModel()).Error; err != nil {
		return nil, err
	}
	return dto, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	// Refuse to delete a category that still has items attached so
	// admins get a clear error instead of orphaning rows.
	var count int64
	if err := r.db.Model(&model.Item{}).
		Where("category_id = ?", id).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return gorm.ErrForeignKeyViolated
	}
	return r.db.Where("id = ?", id).Delete(&model.ItemCategory{}).Error
}

func (r *repository) FindAll(
	ctx context.Context,
	req *entity.ItemCategoryFindAllRequest,
) (*pagination.ResultPagination, error) {
	var rows []model.ItemCategory = make([]model.ItemCategory, 0)
	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		// Categories are org-scoped. NULL organization_id rows are
		// treated as global seed data (visible to every org) so the
		// catalog stays populated during early onboarding.
		q := r.db.Model(&model.ItemCategory{})
		if req.FindAllRequest.OrganizationData.ID != "" {
			q = q.Where(
				"organization_id IS NULL OR organization_id = ?",
				req.FindAllRequest.OrganizationData.ID,
			)
		}
		if s := strings.TrimSpace(req.Query); s != "" {
			like := "%" + s + "%"
			q = q.Where("name ILIKE ?", like)
		}
		return q
	}, &pagination.TableRequest{
		Request:       req,
		QueryField:    []string{},
		Data:          &rows,
		AllowedFields: []string{"name"},
	})
	if err != nil {
		return nil, err
	}
	result := dataTable.(*pagination.ResultPagination)
	hits := result.Data.(*[]model.ItemCategory)
	out := make([]*entity.ItemCategoryDto, 0, len(*hits))
	for i := range *hits {
		out = append(out, entity.NewItemCategoryDtoFromModel(&(*hits)[i]))
	}
	return &pagination.ResultPagination{
		Data:        out,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}
