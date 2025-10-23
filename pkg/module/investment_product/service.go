package investmentproduct

import (
	"context"
	"time"

	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/investment"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type Service interface {
	Create(ctx context.Context, dto *entity.InvestmentProductDto) (*entity.InvestmentProductDto, error)
	FindByID(ctx context.Context, id string, includeAum bool) (*entity.InvestmentProductDto, error)
	Update(ctx context.Context, dto *entity.InvestmentProductDto) (*entity.InvestmentProductDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.InvestmentProductFilter) (*pagination.ResultPagination, error)

	SummaryList(ctx context.Context) ([]*entity.InvestmentProductSummaryDto, error)
	SetCallbackNetAssetValue(cb func(ctx context.Context, investmentProductID string, date string) (*entity.NetAssetValueDto, error))
	CalculateAUM(ctx context.Context, investmentProductID string) (float64, error)
}

type service struct {
	repo            Repository
	cbNetAssetValue func(ctx context.Context, investmentProductID string, date string) (*entity.NetAssetValueDto, error)
}

func NewService(
	repo Repository,
) Service {
	return &service{repo: repo}
}

func (s *service) SetCallbackNetAssetValue(cb func(ctx context.Context, investmentProductID string, date string) (*entity.NetAssetValueDto, error)) {
	s.cbNetAssetValue = cb
}

func (s *service) Create(ctx context.Context, dto *entity.InvestmentProductDto) (*entity.InvestmentProductDto, error) {
	dto.OrganizationID = shared.GetOrganization(ctx).ID
	return s.repo.Create(ctx, dto)
}

func (s *service) FindByID(ctx context.Context, id string, includeAum bool) (*entity.InvestmentProductDto, error) {
	dto, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	dateNow := time.Now().Format("2006-01-02")
	if s.cbNetAssetValue != nil {
		netAssetValue, err := s.cbNetAssetValue(ctx, dto.ID, dateNow)
		if err == nil {
			dto.NetAssetValueDto = netAssetValue
		}
	}

	if includeAum {
		aum, err := s.CalculateAUM(ctx, dto.ID)
		if err == nil {
			dto.AUM = aum
		}
	}

	return dto, nil
}

func (s *service) Update(ctx context.Context, dto *entity.InvestmentProductDto) (*entity.InvestmentProductDto, error) {
	dto.OrganizationID = shared.GetOrganization(ctx).ID
	return s.repo.Update(ctx, dto)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) FindAll(ctx context.Context, req *entity.InvestmentProductFilter) (*pagination.ResultPagination, error) {
	result, err := s.repo.FindAll(ctx, req)
	if err != nil {
		return nil, err
	}

	dateNow := time.Now().Format("2006-01-02")
	data := result.Data.([]*entity.InvestmentProductDto)

	for _, dto := range data {
		if s.cbNetAssetValue != nil {
			netAssetValue, err := s.cbNetAssetValue(ctx, dto.ID, dateNow)
			if err == nil {
				dto.NetAssetValueDto = netAssetValue
			}
		}

		if req.IncludeAum {
			aum, err := s.CalculateAUM(ctx, dto.ID)
			if err == nil {
				dto.AUM = aum
			}
		}
	}

	return result, nil
}

func (s *service) SummaryList(ctx context.Context) ([]*entity.InvestmentProductSummaryDto, error) {
	return s.repo.SummaryList(ctx)
}

func (s *service) CalculateAUM(ctx context.Context, investmentProductID string) (float64, error) {
	return s.repo.CalculateAUM(ctx, investmentProductID)
}
