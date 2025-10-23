package article

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, dto *entity.ArticleDto) (*entity.ArticleDto, error)
	FindByID(ctx context.Context, id string) (*entity.ArticleDto, error)
	Update(ctx context.Context, dto *entity.ArticleDto) (*entity.ArticleDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.ArticleFindAllRequest) (*pagination.ResultPagination, error)
	FindBySlug(ctx context.Context, slug string) (*entity.ArticleDto, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, dto *entity.ArticleDto) (*entity.ArticleDto, error) {
	m := dto.ToModel()
	err := r.db.Create(m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.ArticleDto).FromModel(m), nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*entity.ArticleDto, error) {
	var m *model.Article
	if err := r.db.Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return new(entity.ArticleDto).FromModel(m), nil
}

func (r *repository) Update(ctx context.Context, dto *entity.ArticleDto) (*entity.ArticleDto, error) {
	m := dto.ToModel()
	err := r.db.Updates(m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.ArticleDto).FromModel(m), nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	err := r.db.Where("id = ?", id).Delete(&model.Article{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) FindAll(ctx context.Context, req *entity.ArticleFindAllRequest) (*pagination.ResultPagination, error) {
	var m []model.Article = make([]model.Article, 0)

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		return r.db.Model(&model.Article{})
	}, &pagination.TableRequest{
		Request:       req,
		QueryField:    []string{"title", "content"},
		Data:          &m,
		AllowedFields: []string{"status"},
	})
	if err != nil {
		return nil, err
	}

	result := dataTable.(*pagination.ResultPagination)
	results := result.Data.(*[]model.Article)
	var data []*entity.ArticleDto = make([]*entity.ArticleDto, 0)
	for _, m := range *results {
		data = append(data, new(entity.ArticleDto).FromModel(&m))
	}
	result.Data = data
	return result, nil
}

func (r *repository) FindBySlug(ctx context.Context, slug string) (*entity.ArticleDto, error) {
	var m *model.Article
	if err := r.db.Where("slug = ?", slug).First(&m).Error; err != nil {
		return nil, err
	}
	return new(entity.ArticleDto).FromModel(m), nil
}
