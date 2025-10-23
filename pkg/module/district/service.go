package district

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type Service interface {
	FindAll(ctx context.Context, req *entity.DistrictFindAllRequest) (*pagination.ResultPagination, error)
	FindByID(ctx context.Context, id string) (*entity.DistrictDto, error)
	FindByCode(ctx context.Context, code string) (*entity.DistrictDto, error)
	FindByRegencyID(ctx context.Context, regencyID string) ([]*entity.DistrictDto, error)
	FindByRegencyCode(ctx context.Context, regencyCode string) ([]*entity.DistrictDto, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) FindAll(ctx context.Context, req *entity.DistrictFindAllRequest) (*pagination.ResultPagination, error) {
	return s.repository.FindAll(ctx, req)
}

func (s *service) FindByID(ctx context.Context, id string) (*entity.DistrictDto, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *service) FindByCode(ctx context.Context, code string) (*entity.DistrictDto, error) {
	return s.repository.FindByCode(ctx, code)
}

func (s *service) FindByRegencyID(ctx context.Context, regencyID string) ([]*entity.DistrictDto, error) {
	return s.repository.FindByRegencyID(ctx, regencyID)
}

func (s *service) FindByRegencyCode(ctx context.Context, regencyCode string) ([]*entity.DistrictDto, error) {
	return s.repository.FindByRegencyCode(ctx, regencyCode)
}
