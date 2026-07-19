package company

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type Service interface {
	FindCompanyByUserID(ctx context.Context, userID string) (*entity.CompanyDto, error)
	FindCompanyByAdminID(ctx context.Context, adminID string) (*entity.CompanyDto, error)
	FindAll(ctx context.Context, req *entity.CompanyFindAllRequest) (*pagination.ResultPagination, error)
}

type service struct {
	repository Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) FindCompanyByUserID(ctx context.Context, userID string) (*entity.CompanyDto, error) {
	return s.repository.FindCompanyByUserID(ctx, userID)
}

func (s *service) FindCompanyByAdminID(ctx context.Context, adminID string) (*entity.CompanyDto, error) {
	return s.repository.FindCompanyByAdminID(ctx, adminID)
}

// FindAll mirrors the org-scope pattern used by item_category and
// salary_component: fill the org id from the request context if
// the caller did not pin one, then let the repository apply the
// WHERE clause.
func (s *service) FindAll(
	ctx context.Context,
	req *entity.CompanyFindAllRequest,
) (*pagination.ResultPagination, error) {
	if req.FindAllRequest.OrganizationData.ID == "" {
		req.FindAllRequest.OrganizationData.ID = shared.GetOrganization(ctx).ID
	}
	return s.repository.FindAll(ctx, req)
}
