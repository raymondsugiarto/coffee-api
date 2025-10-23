package userlog

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	e "github.com/raymondsugiarto/coffee-api/pkg/entity"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type Service interface {
	Create(ctx context.Context, dto *entity.UserLogDto) (*entity.UserLogDto, error)
	FindByID(ctx context.Context, id string) (*entity.UserLogDto, error)
	Update(ctx context.Context, dto *entity.UserLogDto) (*entity.UserLogDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *e.UserLogFindAllRequest) (*pagination.ResultPagination, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Create(ctx context.Context, dto *entity.UserLogDto) (*entity.UserLogDto, error) {
	dto.OrganizationID = shared.GetOrganization(ctx).ID
	return s.repo.Create(ctx, dto)
}

func (s *service) FindByID(ctx context.Context, id string) (*entity.UserLogDto, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) Update(ctx context.Context, dto *entity.UserLogDto) (*entity.UserLogDto, error) {
	dto.OrganizationID = shared.GetOrganization(ctx).ID
	return s.repo.Update(ctx, dto)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) FindAll(ctx context.Context, req *e.UserLogFindAllRequest) (*pagination.ResultPagination, error) {
	return s.repo.FindAll(ctx, req)
}
