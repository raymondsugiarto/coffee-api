package portfolio

import (
	"context"

	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/report"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type Service interface {
	FindAllPortfolioWithNav(ctx context.Context, req *entity.PortfolioFindAllRequest) (*pagination.ResultPagination, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) FindAllPortfolioWithNav(ctx context.Context, req *entity.PortfolioFindAllRequest) (*pagination.ResultPagination, error) {
	return s.repository.FindAllPortfolioWithNav(ctx, req)
}
