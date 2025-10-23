package unitlink

import (
	"context"
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	ce "github.com/raymondsugiarto/coffee-api/pkg/entity/customer"
	investmentitem "github.com/raymondsugiarto/coffee-api/pkg/module/investment/investment_item"
	investmentproduct "github.com/raymondsugiarto/coffee-api/pkg/module/investment_product"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

const (
	defaultSize   = -1
	defaultAmount = 0
)

type Service interface {
	Create(ctx context.Context, dto *ce.UnitLinkDto) (*ce.UnitLinkDto, error)
	CreateWithTx(ctx context.Context, tx *gorm.DB, dto *ce.UnitLinkDto) (*ce.UnitLinkDto, error)
	Get(ctx context.Context, id string) (*ce.UnitLinkDto, error)
	Update(ctx context.Context, dto *ce.UnitLinkDto) (*ce.UnitLinkDto, error)
	UpdateNab(ctx context.Context, dto *ce.UnitLinkDto) (*ce.UnitLinkDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *ce.UnitLinkFindAllRequest) (*pagination.ResultPagination, error)
	FindAllByTransactionDate(ctx context.Context, transactionDate time.Time) ([]*ce.UnitLinkDto, error)
	FindAllInvestmentProductByParticipant(ctx context.Context, participantID string) ([]*ce.UnitLinkPortfolioDto, error)
	FindAllInvestmentProductByCustomer(ctx context.Context, customerID string) ([]*ce.UnitLinkPortfolioDto, error)
	SumInvestmentProductByCustomer(ctx context.Context, customerID string) (*ce.SumUnitLinkPortfolioDto, error)
	SumInvestmentProductByParticipant(ctx context.Context, participantID string) (*ce.SumUnitLinkPortfolioDto, error)
	ClaimUnitLink(ctx context.Context, tx *gorm.DB, participantID string) error
	FindAllInvestmentProductGroupParticipant(ctx context.Context) ([]*ce.UnitLinkPortfolioGroupParticipantDto, error)
	FindLatestEachProductAndParticipantAndType(ctx context.Context) ([]*ce.UnitLinkLatestEachProductAndParticipantAndTypeDto, error)
	SummaryByCompany(ctx context.Context) (*ce.UnitLinkSummaryCompanyDto, error)
	SummaryPerType(ctx context.Context) ([]*ce.UnitLinkSummaryPerTypeDto, error)
	FindAllPortfolioWithNav(ctx context.Context, req *ce.PortfolioFindAllRequest) (*pagination.ResultPagination, error)
}

type service struct {
	repo                 Repository
	investmentItemSvc    investmentitem.Service
	investmentProductSvc investmentproduct.Service
}

func NewService(repo Repository, investmentItemSvc investmentitem.Service, investmentProductSvc investmentproduct.Service) Service {
	return &service{repo, investmentItemSvc, investmentProductSvc}
}

func (s *service) CreateWithTx(ctx context.Context, tx *gorm.DB, dto *ce.UnitLinkDto) (*ce.UnitLinkDto, error) {
	m, err := s.repo.CreateWithTx(ctx, tx, dto)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *service) Create(ctx context.Context, dto *ce.UnitLinkDto) (*ce.UnitLinkDto, error) {
	m, err := s.repo.Create(ctx, dto)
	if err != nil {
		return nil, err
	}
	return m, nil
}
func (s *service) Get(ctx context.Context, id string) (*ce.UnitLinkDto, error) {
	m, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *service) Update(ctx context.Context, dto *ce.UnitLinkDto) (*ce.UnitLinkDto, error) {
	m, err := s.repo.Update(ctx, dto)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *service) UpdateNab(ctx context.Context, dto *ce.UnitLinkDto) (*ce.UnitLinkDto, error) {
	m, err := s.repo.UpdateNab(ctx, dto)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) FindAll(ctx context.Context, req *ce.UnitLinkFindAllRequest) (*pagination.ResultPagination, error) {
	m, err := s.repo.FindAll(ctx, req)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *service) FindAllByTransactionDate(ctx context.Context, transactionDate time.Time) ([]*ce.UnitLinkDto, error) {
	m, err := s.repo.FindAllByTransactionDate(ctx, transactionDate)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *service) FindAllInvestmentProductByCustomer(ctx context.Context, customerID string) ([]*ce.UnitLinkPortfolioDto, error) {
	m, err := s.repo.FindAllInvestmentProductByCustomer(ctx, customerID)
	if err != nil {
		return nil, err
	}

	// Load investment products separately for each portfolio item
	for _, portfolio := range m {
		investmentProduct, err := s.investmentProductSvc.FindByID(ctx, portfolio.InvestmentProductID, false)
		if err != nil {
			return nil, err
		}
		portfolio.InvestmentProduct = *investmentProduct
	}

	return m, nil
}

func (s *service) FindAllInvestmentProductByParticipant(ctx context.Context, participantID string) ([]*ce.UnitLinkPortfolioDto, error) {
	m, err := s.repo.FindAllInvestmentProductByParticipant(ctx, participantID)
	if err != nil {
		return nil, err
	}

	// Load investment products separately for each portfolio item
	for _, portfolio := range m {
		investmentProduct, err := s.investmentProductSvc.FindByID(ctx, portfolio.InvestmentProductID, false)
		if err != nil {
			return nil, err
		}
		portfolio.InvestmentProduct = *investmentProduct
	}

	return m, nil
}

func (s *service) SumInvestmentProductByCustomer(ctx context.Context, customerID string) (*ce.SumUnitLinkPortfolioDto, error) {
	m, err := s.repo.SumInvestmentProductByCustomer(ctx, customerID)
	if err != nil {
		return nil, err
	}

	m.Profit = m.CurrentBalance - m.TotalModal

	if m.TotalModal > 0 {
		m.ReturnPercentage = (m.Profit / m.TotalModal) * 100
	} else {
		m.ReturnPercentage = 0
	}

	return m, nil
}

func (s *service) SumInvestmentProductByParticipant(ctx context.Context, participantID string) (*ce.SumUnitLinkPortfolioDto, error) {
	m, err := s.repo.SumInvestmentProductByParticipant(ctx, participantID)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *service) FindAllInvestmentProductGroupParticipant(ctx context.Context) ([]*ce.UnitLinkPortfolioGroupParticipantDto, error) {
	m, err := s.repo.FindAllInvestmentProductGroupParticipant(ctx)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (s *service) FindLatestEachProductAndParticipantAndType(ctx context.Context) ([]*ce.UnitLinkLatestEachProductAndParticipantAndTypeDto, error) {
	m, err := s.repo.FindLatestEachProductAndParticipantAndType(ctx)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (s *service) ClaimUnitLink(ctx context.Context, tx *gorm.DB, participantID string) error {
	unitLinkPage, err := s.repo.FindAll(ctx, &ce.UnitLinkFindAllRequest{
		FindAllRequest: entity.FindAllRequest{
			GetListRequest: pagination.GetListRequest{
				Size: defaultSize,
			},
		},
		ParticipantID: participantID,
	})
	if err != nil {
		return err
	}

	unitLinks := unitLinkPage.Data.([]*ce.UnitLinkDto)

	for _, item := range unitLinks {
		item.TotalAmount = defaultAmount
		item.Nab = defaultAmount
		item.Ip = defaultAmount

		_, err := s.repo.UpdateTotalAmountNabIpWithTx(ctx, tx, item)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *service) SummaryByCompany(ctx context.Context) (*ce.UnitLinkSummaryCompanyDto, error) {
	m, err := s.repo.SummaryByCompany(ctx, *shared.GetCompanyID(ctx))
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (s *service) SummaryPerType(ctx context.Context) ([]*ce.UnitLinkSummaryPerTypeDto, error) {
	m, err := s.repo.SummaryPerType(ctx)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (s *service) FindAllPortfolioWithNav(ctx context.Context, req *ce.PortfolioFindAllRequest) (*pagination.ResultPagination, error) {
	m, err := s.repo.FindAllPortfolioWithNav(ctx, req)
	if err != nil {
		return nil, err
	}
	return m, nil
}
