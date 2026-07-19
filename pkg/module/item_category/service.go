package itemcategory

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type Service interface {
	Create(ctx context.Context, dto *entity.ItemCategoryDto) (*entity.ItemCategoryDto, error)
	Get(ctx context.Context, id string) (*entity.ItemCategoryDto, error)
	Update(ctx context.Context, dto *entity.ItemCategoryDto) (*entity.ItemCategoryDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.ItemCategoryFindAllRequest) (*pagination.ResultPagination, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, dto *entity.ItemCategoryDto) (*entity.ItemCategoryDto, error) {
	if dto.OrganizationID == "" {
		dto.OrganizationID = shared.GetOrganization(ctx).ID
	}
	return s.repo.Create(ctx, dto)
}

func (s *service) Get(ctx context.Context, id string) (*entity.ItemCategoryDto, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) Update(ctx context.Context, dto *entity.ItemCategoryDto) (*entity.ItemCategoryDto, error) {
	if dto.OrganizationID == "" {
		dto.OrganizationID = shared.GetOrganization(ctx).ID
	}
	return s.repo.Update(ctx, dto)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) FindAll(
	ctx context.Context,
	req *entity.ItemCategoryFindAllRequest,
) (*pagination.ResultPagination, error) {
	// Categories are org-scoped. Mirror the stock_session pattern:
	// fill the org id from the request context if the caller did
	// not already pin one, then let the repository apply WHERE.
	if req.FindAllRequest.OrganizationData.ID == "" {
		req.FindAllRequest.OrganizationData.ID = shared.GetOrganization(ctx).ID
	}
	return s.repo.FindAll(ctx, req)
}
