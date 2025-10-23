package article

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type Service interface {
	Create(ctx context.Context, dto *entity.ArticleDto) (*entity.ArticleDto, error)
	FindByID(ctx context.Context, id string) (*entity.ArticleDto, error)
	Update(ctx context.Context, dto *entity.ArticleDto) (*entity.ArticleDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.ArticleFindAllRequest) (*pagination.ResultPagination, error)
	FindBySlug(ctx context.Context, slug string) (*entity.ArticleDto, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository}
}

func (s *service) Create(ctx context.Context, dto *entity.ArticleDto) (*entity.ArticleDto, error) {
	userCredential := shared.GetUserCredential(ctx)
	dto.CreatedBy = userCredential.UserID
	dto.UpdatedBy = userCredential.UserID
	return s.repository.Create(ctx, dto)
}

func (s *service) FindByID(ctx context.Context, id string) (*entity.ArticleDto, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *service) Update(ctx context.Context, dto *entity.ArticleDto) (*entity.ArticleDto, error) {
	userCredential := shared.GetUserCredential(ctx)
	dto.UpdatedBy = userCredential.UserID
	return s.repository.Update(ctx, dto)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}

func (s *service) FindAll(ctx context.Context, req *entity.ArticleFindAllRequest) (*pagination.ResultPagination, error) {
	userType := shared.GetOriginTypeKey(ctx)
	if userType != string(entity.ADMIN) {
		req.Status = entity.PUBLISHED
	}
	return s.repository.FindAll(ctx, req)
}

func (s *service) FindBySlug(ctx context.Context, slug string) (*entity.ArticleDto, error) {
	return s.repository.FindBySlug(ctx, slug)
}
