package accounting

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

// AccountMutationRepository is the ledger-side persistence
// contract. Subset of the full CRUD surface because the ledger
// is append-only.
type AccountMutationRepository interface {
	Create(ctx context.Context, dto *entity.AccountMutationDto) (*entity.AccountMutationDto, error)
	FindAll(ctx context.Context, req *entity.AccountMutationFindAllRequest) (*pagination.ResultPagination, error)
}

// AccountMutationService is the wired-side contract. Exposed to
// upstream flows (stock_session close, order paid, payroll
// posted, ...) so they can post to the ledger without needing to
// know the SQL.
type AccountMutationService interface {
	Create(ctx context.Context, dto *entity.AccountMutationDto) (*entity.AccountMutationDto, error)
	// Post writes the same payload but, in addition to the
	// create, runs the account-existence check against the
	// account repo. It returns an error if account_id doesn't
	// resolve, so callers can fail-fast before persisting a bad
	// pointer. This is the canonical entry point for any flow
	// posting a mutation: write-through validation > implicit FK
	// errors at insert time.
	Post(ctx context.Context, dto *entity.AccountMutationDto) (*entity.AccountMutationDto, error)
	FindAll(ctx context.Context, req *entity.AccountMutationFindAllRequest) (*pagination.ResultPagination, error)
}

type accountMutationService struct {
	repo            AccountMutationRepository
	accountResolver AccountResolver
}

// AccountResolver is a narrow interface the mutation service
// uses to validate that account_id points at a real row before
// posting. The full AccountService satisfies this without us
// having to depend on it directly — keeping this service test-
// friendly and free of graph-store coupling.
type AccountResolver interface {
	Get(ctx context.Context, id string) (*entity.AccountDto, error)
}

// NewAccountMutationService wires the mutation ledger. Pass the
// full AccountService in — it satisfies AccountResolver via its
// Get method.
func NewAccountMutationService(
	repo AccountMutationRepository,
	accountResolver AccountResolver,
) AccountMutationService {
	return &accountMutationService{
		repo:            repo,
		accountResolver: accountResolver,
	}
}

func (s *accountMutationService) Create(
	ctx context.Context,
	dto *entity.AccountMutationDto,
) (*entity.AccountMutationDto, error) {
	if dto.OrganizationID == "" {
		dto.OrganizationID = shared.GetOrganization(ctx).ID
	}
	return s.repo.Create(ctx, dto)
}

func (s *accountMutationService) Post(
	ctx context.Context,
	dto *entity.AccountMutationDto,
) (*entity.AccountMutationDto, error) {
	// Refuse to post a mutation that points at a nonexistent
	// account. Returning a 400-ish error from here surfaces a
	// clean validation message instead of letting the FK fire
	// later in the SQL.
	if _, err := s.accountResolver.Get(ctx, dto.AccountID); err != nil {
		return nil, err
	}
	return s.Create(ctx, dto)
}

func (s *accountMutationService) FindAll(
	ctx context.Context,
	req *entity.AccountMutationFindAllRequest,
) (*pagination.ResultPagination, error) {
	if req.FindAllRequest.OrganizationData.ID == "" {
		req.FindAllRequest.OrganizationData.ID = shared.GetOrganization(ctx).ID
	}
	return s.repo.FindAll(ctx, req)
}
