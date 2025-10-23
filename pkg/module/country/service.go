package country

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type Service interface {
	FindAll(ctx context.Context, req *entity.CountryFindAllRequest) (*pagination.ResultPagination, error)
	FindByID(ctx context.Context, id string) (*entity.CountryDto, error)
	FindByCCA2(ctx context.Context, cca2 string) (*entity.CountryDto, error)
	FindByCCA3(ctx context.Context, cca3 string) (*entity.CountryDto, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) FindAll(ctx context.Context, req *entity.CountryFindAllRequest) (*pagination.ResultPagination, error) {
	return s.repository.FindAll(ctx, req)
}

func (s *service) FindByID(ctx context.Context, id string) (*entity.CountryDto, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *service) FindByCCA2(ctx context.Context, cca2 string) (*entity.CountryDto, error) {
	return s.repository.FindByCCA2(ctx, cca2)
}

func (s *service) FindByCCA3(ctx context.Context, cca3 string) (*entity.CountryDto, error) {
	return s.repository.FindByCCA3(ctx, cca3)
}
