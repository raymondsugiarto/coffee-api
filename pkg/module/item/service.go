package item

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type Service interface {
	Create(ctx context.Context, dto *entity.ItemDto) (*entity.ItemDto, error)
	FindByID(ctx context.Context, id string) (*entity.ItemDto, error)
	Update(ctx context.Context, dto *entity.ItemDto) (*entity.ItemDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.ItemFindAllRequest) (*pagination.ResultPagination, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, dto *entity.ItemDto) (*entity.ItemDto, error) {
	dto.OrganizationID = shared.GetOrganization(ctx).ID
	return s.repo.Create(ctx, dto)
}

func (s *service) FindByID(ctx context.Context, id string) (*entity.ItemDto, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) Update(ctx context.Context, dto *entity.ItemDto) (*entity.ItemDto, error) {
	dto.OrganizationID = shared.GetOrganization(ctx).ID
	return s.repo.Update(ctx, dto)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) FindAll(ctx context.Context, req *entity.ItemFindAllRequest) (*pagination.ResultPagination, error) {
	// Items are now scoped to the active organization only — there
	// is no longer a `company_id` join in the find-all query. Mirror
	// the stock_session pattern: pull the org id off the request
	// context and let the repository apply the WHERE clause.
	if req.FindAllRequest.OrganizationData.ID == "" {
		req.FindAllRequest.OrganizationData.ID = shared.GetOrganization(ctx).ID
	}
	return s.repo.FindAll(ctx, req)
}
