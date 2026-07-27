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
	// FindByCompany returns the salary bands for one company,
	// ordered by minimum_target ASC. Used by the stock-session
	// close path so the picking logic can walk the list directly.
	FindByCompany(ctx context.Context, companyID string) ([]*entity.SalaryComponentDto, error)
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

// FindByCompany is the canonical entry point that the
// stock-session close path uses to resolve salary bands. We pass
// the company id through unmodified — the repository handles the
// empty-company short-circuit so the caller doesn't need to.
//
// Errors are propagated as-is: callers (e.g. stock_session) treat
// a non-nil error as "salary lookup failed, fall back to a 0-amount
// breakdown" rather than aborting the close, so we don't wrap.
func (s *service) FindByCompany(
	ctx context.Context,
	companyID string,
) ([]*entity.SalaryComponentDto, error) {
	return s.repo.FindByCompany(ctx, companyID)
}
