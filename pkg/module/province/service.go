package province

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type Service interface {
	FindAll(ctx context.Context, req *entity.ProvinceFindAllRequest) (*pagination.ResultPagination, error)
	FindByID(ctx context.Context, id string) (*entity.ProvinceDto, error)
	FindByCode(ctx context.Context, code string) (*entity.ProvinceDto, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) FindAll(ctx context.Context, req *entity.ProvinceFindAllRequest) (*pagination.ResultPagination, error) {
	return s.repository.FindAll(ctx, req)
}

func (s *service) FindByID(ctx context.Context, id string) (*entity.ProvinceDto, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *service) FindByCode(ctx context.Context, code string) (*entity.ProvinceDto, error) {
	return s.repository.FindByCode(ctx, code)
}
