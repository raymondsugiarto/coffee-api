package permission

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type Service interface {
	FindAll(ctx context.Context, req *entity.PermissionFindAllRequest) (*pagination.ResultPagination, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) FindAll(ctx context.Context, req *entity.PermissionFindAllRequest) (*pagination.ResultPagination, error) {
	return s.repo.FindAll(ctx, req)
}
