package village

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type Service interface {
	FindAll(ctx context.Context, req *entity.VillageFindAllRequest) (*pagination.ResultPagination, error)
	FindByID(ctx context.Context, id string) (*entity.VillageDto, error)
	FindByCode(ctx context.Context, code string) (*entity.VillageDto, error)
	FindByDistrictID(ctx context.Context, districtId string) ([]*entity.VillageDto, error)
	FindByDistrictCode(ctx context.Context, districtCode string) ([]*entity.VillageDto, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) FindAll(ctx context.Context, req *entity.VillageFindAllRequest) (*pagination.ResultPagination, error) {
	return s.repository.FindAll(ctx, req)
}

func (s *service) FindByID(ctx context.Context, id string) (*entity.VillageDto, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *service) FindByCode(ctx context.Context, code string) (*entity.VillageDto, error) {
	return s.repository.FindByCode(ctx, code)
}

func (s *service) FindByDistrictID(ctx context.Context, districtId string) ([]*entity.VillageDto, error) {
	return s.repository.FindByDistrictID(ctx, districtId)
}

func (s *service) FindByDistrictCode(ctx context.Context, districtCode string) ([]*entity.VillageDto, error) {
	return s.repository.FindByDistrictCode(ctx, districtCode)
}
