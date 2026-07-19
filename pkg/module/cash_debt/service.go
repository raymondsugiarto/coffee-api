package cashdebt

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type Service interface {
	Create(ctx context.Context, dto *entity.CashDebtDto) (*entity.CashDebtDto, error)
	Get(ctx context.Context, id string) (*entity.CashDebtDto, error)
	Update(ctx context.Context, dto *entity.CashDebtDto) (*entity.CashDebtDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.CashDebtFindAllRequest) (*pagination.ResultPagination, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, dto *entity.CashDebtDto) (*entity.CashDebtDto, error) {
	if dto.OrganizationID == "" {
		dto.OrganizationID = shared.GetOrganization(ctx).ID
	}
	return s.repo.Create(ctx, dto)
}

func (s *service) Get(ctx context.Context, id string) (*entity.CashDebtDto, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) Update(ctx context.Context, dto *entity.CashDebtDto) (*entity.CashDebtDto, error) {
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
	req *entity.CashDebtFindAllRequest,
) (*pagination.ResultPagination, error) {
	if req.FindAllRequest.OrganizationData.ID == "" {
		req.FindAllRequest.OrganizationData.ID = shared.GetOrganization(ctx).ID
	}
	return s.repo.FindAll(ctx, req)
}
