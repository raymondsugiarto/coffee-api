package investmentitem

import (
	"context"
	"errors"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	customer "github.com/raymondsugiarto/coffee-api/pkg/entity/customer"
	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/investment"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Service interface {
	Create(ctx context.Context, dto *entity.InvestmentItemDto) (*entity.InvestmentItemDto, error)
	FindByID(ctx context.Context, id string) (*entity.InvestmentItemDto, error)
	Update(ctx context.Context, dto *entity.InvestmentItemDto) (*entity.InvestmentItemDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.InvestmentItemFindAllRequest) (*pagination.ResultPagination, error)

	CreateBatchWithTx(ctx context.Context, db *gorm.DB, dto []*entity.InvestmentItemDto) ([]*entity.InvestmentItemDto, error)
	CreateWithTxAndCallback(ctx context.Context, db *gorm.DB, dto *entity.InvestmentDto) error
	SumInvestmentByCustomer(ctx context.Context, customerID string) (*customer.SumUnitLinkPortfolioDto, error)
	FindByInvestmentID(ctx context.Context, investmentID string) ([]*entity.InvestmentItemDto, error)
	HasPreviousInvestmentForCustomer(ctx context.Context, investmentID string, customerID string) (bool, error)
	PrepareForStatement(ctx context.Context, req *entity.InvestmentItemFindAllRequest) ([]*entity.InvestmentStatementDto, error)
}

type service struct {
	repo Repository
}

func NewService(
	repo Repository,
) Service {
	return &service{repo}
}

func (s *service) CreateWithTxAndCallback(ctx context.Context, db *gorm.DB, dto *entity.InvestmentDto) error {
	_, err := s.CreateBatchWithTx(ctx, db, dto.InvestmentItems)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) CreateBatchWithTx(ctx context.Context, db *gorm.DB, dto []*entity.InvestmentItemDto) ([]*entity.InvestmentItemDto, error) {
	// Generate codes for investment items that don't have them
	for _, item := range dto {
		if item.Code == "" {
			code, err := s.generateUniqueInvestmentItemCode(ctx)
			if err != nil {
				return nil, err
			}
			item.Code = code
		}
	}
	return s.repo.CreateBatchWithTx(ctx, db, dto)
}

func (s *service) Create(ctx context.Context, dto *entity.InvestmentItemDto) (*entity.InvestmentItemDto, error) {
	dto.OrganizationID = shared.GetOrganization(ctx).ID
	return s.repo.Create(ctx, dto)
}

func (s *service) FindByID(ctx context.Context, id string) (*entity.InvestmentItemDto, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) Update(ctx context.Context, dto *entity.InvestmentItemDto) (*entity.InvestmentItemDto, error) {
	dto.OrganizationID = shared.GetOrganization(ctx).ID
	return s.repo.Update(ctx, dto)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) FindAll(ctx context.Context, req *entity.InvestmentItemFindAllRequest) (*pagination.ResultPagination, error) {
	return s.repo.FindAll(ctx, req)
}

func (s *service) SumInvestmentByCustomer(ctx context.Context, customerID string) (*customer.SumUnitLinkPortfolioDto, error) {
	return s.repo.SumInvestmentByCustomer(ctx, customerID)
}

func (s *service) FindByInvestmentID(ctx context.Context, investmentID string) ([]*entity.InvestmentItemDto, error) {
	return s.repo.FindByInvestmentID(ctx, investmentID)
}

func (s *service) HasPreviousInvestmentForCustomer(ctx context.Context, investmentID string, customerID string) (bool, error) {
	return s.repo.HasPreviousInvestmentForCustomer(ctx, investmentID, customerID)
}

func (s *service) PrepareForStatement(ctx context.Context, req *entity.InvestmentItemFindAllRequest) ([]*entity.InvestmentStatementDto, error) {
	itemPage, err := s.FindAll(ctx, req)
	if err != nil {
		return nil, err
	}

	items := itemPage.Data.([]*entity.InvestmentItemDto)

	statementItems := make([]*entity.InvestmentStatementDto, len(items))
	for i, item := range items {
		statementItems[i] = item.ToInvestmentStatementDto()
	}

	err = s.enhanceWithNavCalculations(ctx, statementItems)
	if err != nil {
		return nil, err
	}

	return statementItems, nil
}

func (s *service) enhanceWithNavCalculations(ctx context.Context, items []*entity.InvestmentStatementDto) error {
	if len(items) == 0 {
		return nil
	}

	productIDs := s.extractUniqueProductIDs(items)

	// Get latest NAV untuk current value
	currentDate := time.Now()
	latestNavMap, err := s.repo.GetNavByProductIDsWithDate(ctx, productIDs, currentDate.Format("2006-01-02"))
	if err != nil {
		return err
	}

	// Process items
	for _, item := range items {
		if navWithDate, exists := latestNavMap[item.InvestmentProductID]; exists {
			item.CurrentNavAmount = navWithDate.Amount
			item.CurrentNavDate = navWithDate.Date

			// Handle unit = 0 case: keep unit as 0 and preserve total amount as current value
			if item.Unit == 0 {
				item.CurrentValue = item.TotalAmount
				item.GainLoss = 0
				item.GainLossPercentage = 0
			} else {
				item.CurrentValue = item.Unit * navWithDate.Amount
				item.GainLoss = item.CurrentValue - item.TotalAmount

				if item.TotalAmount > 0 {
					item.GainLossPercentage = (item.GainLoss / item.TotalAmount) * 100
				}
			}
		}
	}

	return nil
}

func (s *service) extractUniqueProductIDs(items []*entity.InvestmentStatementDto) []string {
	productIDMap := make(map[string]bool)
	for _, item := range items {
		productIDMap[item.InvestmentProductID] = true
	}

	productIDs := make([]string, 0, len(productIDMap))
	for productID := range productIDMap {
		productIDs = append(productIDs, productID)
	}

	return productIDs
}

func (s *service) generateUniqueInvestmentItemCode(ctx context.Context) (string, error) {
	const maxRetries = 10

	for range maxRetries {
		// Generate 7 digit number (0000000 - 9999999)
		code, err := gonanoid.Generate("0123456789", 7)
		if err != nil {
			return "", err
		}

		_, err = s.repo.GetByCode(ctx, code)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return code, nil
			}
			return "", err
		}
	}

	return "", errors.New("failed to generate unique investment item code after maximum retries")
}
