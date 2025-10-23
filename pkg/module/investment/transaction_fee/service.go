package transactionfee

import (
	"context"
	"time"

	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/investment"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
)

type Service interface {
	Create(ctx context.Context, dto *entity.TransactionFeeDto) (*entity.TransactionFeeDto, error)
	GetOperationalFeesByCompanyAndDateRange(ctx context.Context, companyID string, startDate, endDate time.Time) (map[string]float64, error)
}

type service struct {
	repo Repository
}

func NewService(
	repo Repository,
) Service {
	return &service{repo}
}

func (s *service) Create(ctx context.Context, dto *entity.TransactionFeeDto) (*entity.TransactionFeeDto, error) {
	dto.OrganizationID = shared.GetOrganization(ctx).ID
	return s.repo.Create(ctx, dto)
}

func (s *service) GetOperationalFeesByCompanyAndDateRange(ctx context.Context, companyID string, startDate, endDate time.Time) (map[string]float64, error) {
	return s.repo.GetOperationalFeesByCompanyAndDateRange(ctx, companyID, startDate, endDate)
}
