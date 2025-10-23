package regency

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type Service interface {
	FindAll(ctx context.Context, req *entity.RegencyFindAllRequest) (*pagination.ResultPagination, error)
	FindByID(ctx context.Context, id string) (*entity.RegencyDto, error)
	FindByCode(ctx context.Context, code string) (*entity.RegencyDto, error)
	FindByProvinceID(ctx context.Context, provinceID string) ([]*entity.RegencyDto, error)
	FindByProvinceCode(ctx context.Context, provinceCode string) ([]*entity.RegencyDto, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) FindAll(ctx context.Context, req *entity.RegencyFindAllRequest) (*pagination.ResultPagination, error) {
	return s.repository.FindAll(ctx, req)
}

func (s *service) FindByID(ctx context.Context, id string) (*entity.RegencyDto, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *service) FindByCode(ctx context.Context, code string) (*entity.RegencyDto, error) {
	return s.repository.FindByCode(ctx, code)
}

func (s *service) FindByProvinceID(ctx context.Context, provinceID string) ([]*entity.RegencyDto, error) {
	return s.repository.FindByProvinceID(ctx, provinceID)
}

func (s *service) FindByProvinceCode(ctx context.Context, provinceCode string) ([]*entity.RegencyDto, error) {
	return s.repository.FindByProvinceCode(ctx, provinceCode)
}
