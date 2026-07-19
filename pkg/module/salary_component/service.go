package salarycomponent

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type Service interface {
	Create(ctx context.Context, dto *entity.SalaryComponentDto) (*entity.SalaryComponentDto, error)
	Get(ctx context.Context, id string) (*entity.SalaryComponentDto, error)
	Update(ctx context.Context, dto *entity.SalaryComponentDto) (*entity.SalaryComponentDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.SalaryComponentFindAllRequest) (*pagination.ResultPagination, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, dto *entity.SalaryComponentDto) (*entity.SalaryComponentDto, error) {
	if dto.OrganizationID == "" {
		dto.OrganizationID = shared.GetOrganization(ctx).ID
	}
	return s.repo.Create(ctx, dto)
}

func (s *service) Get(ctx context.Context, id string) (*entity.SalaryComponentDto, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) Update(ctx context.Context, dto *entity.SalaryComponentDto) (*entity.SalaryComponentDto, error) {
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
	req *entity.SalaryComponentFindAllRequest,
) (*pagination.ResultPagination, error) {
	// Mirror the org-scope pattern used by item, item_category, and
	// stock_session: fill the org id from the request context if
	// the caller didn't pin one, then let the repository apply the
	// WHERE clause.
	if req.FindAllRequest.OrganizationData.ID == "" {
		req.FindAllRequest.OrganizationData.ID = shared.GetOrganization(ctx).ID
	}
	return s.repo.FindAll(ctx, req)
}
