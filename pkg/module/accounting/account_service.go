package accounting

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

// AccountRepository is the persistence-side contract the service
// depends on. Kept as an interface so tests can substitute an
// in-memory implementation without touching GORM.
type AccountRepository interface {
	Create(ctx context.Context, dto *entity.AccountDto) (*entity.AccountDto, error)
	Get(ctx context.Context, id string) (*entity.AccountDto, error)
	Update(ctx context.Context, dto *entity.AccountDto) (*entity.AccountDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.AccountFindAllRequest) (*pagination.ResultPagination, error)
}

// AccountService is the wired-side interface exposed to handlers
// and other modules (e.g. account_mutation service uses it to
// validate the `account_id` reference at write time).
type AccountService interface {
	Create(ctx context.Context, dto *entity.AccountDto) (*entity.AccountDto, error)
	Get(ctx context.Context, id string) (*entity.AccountDto, error)
	Update(ctx context.Context, dto *entity.AccountDto) (*entity.AccountDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.AccountFindAllRequest) (*pagination.ResultPagination, error)
}

type accountService struct {
	repo AccountRepository
}

func NewAccountService(repo AccountRepository) AccountService {
	return &accountService{repo: repo}
}

func (s *accountService) Create(ctx context.Context, dto *entity.AccountDto) (*entity.AccountDto, error) {
	// Stamp the org id from the request context if the caller
	// didn't set one explicitly. Mirror of the salarycomponent /
	// item_category convention.
	if dto.OrganizationID == "" {
		dto.OrganizationID = shared.GetOrganization(ctx).ID
	}
	return s.repo.Create(ctx, dto)
}

func (s *accountService) Get(ctx context.Context, id string) (*entity.AccountDto, error) {
	return s.repo.Get(ctx, id)
}

func (s *accountService) Update(ctx context.Context, dto *entity.AccountDto) (*entity.AccountDto, error) {
	if dto.OrganizationID == "" {
		dto.OrganizationID = shared.GetOrganization(ctx).ID
	}
	return s.repo.Update(ctx, dto)
}

func (s *accountService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *accountService) FindAll(
	ctx context.Context,
	req *entity.AccountFindAllRequest,
) (*pagination.ResultPagination, error) {
	if req.FindAllRequest.OrganizationData.ID == "" {
		req.FindAllRequest.OrganizationData.ID = shared.GetOrganization(ctx).ID
	}
	return s.repo.FindAll(ctx, req)
}
