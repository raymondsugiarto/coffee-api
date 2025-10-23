package ojkcustomerreport

import (
	"context"
	"fmt"
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	crEntity "github.com/raymondsugiarto/coffee-api/pkg/entity/customer_report"
	ie "github.com/raymondsugiarto/coffee-api/pkg/entity/investment"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"gorm.io/gorm"
)

type Repository interface {
	GetOJKCustomerReportData(ctx context.Context, filter *crEntity.OJKCustomerReportFilterDto) (*OJKCustomerReportData, error)
	GetCustomerBasicInfo(ctx context.Context, customerID string) (*entity.CustomerDto, error)
	GetInvestmentSummary(ctx context.Context, customerID string, startDate, endDate time.Time) (*InvestmentSummaryData, error)
	GetNAVData(ctx context.Context, endDate time.Time) (map[string]*ie.NetAssetValueDto, error)
	GetInvestmentDistributions(ctx context.Context, customerID string) ([]*ie.InvestmentDistributionDto, error)
	GetOpeningBalance(ctx context.Context, customerID string, beforeDate time.Time) (map[string]float64, error)
	GetHistoricalCost(ctx context.Context, customerID string, endDate time.Time) (float64, error)
	GetCurrentPortfolioValue(ctx context.Context, customerID string, endDate time.Time) (float64, error)
	GetOperationalFees(ctx context.Context, customerID string, startDate, endDate time.Time) (map[string]float64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

type OJKCustomerReportData struct {
	Customer                *entity.CustomerDto
	InvestmentItems         []*ie.InvestmentItemDto
	NAVData                 map[string]*ie.NetAssetValueDto
	InvestmentDistributions []*ie.InvestmentDistributionDto
	ProductAmounts          map[string]ProductAmounts
	TotalContribution       float64
	TotalFees               float64
	TotalOperationalFees    float64
	OpeningBalance          map[string]float64
	HistoricalCost          float64
	CurrentPortfolioValue   float64
}

type InvestmentSummaryData struct {
	Items                []*ie.InvestmentItemDto
	TotalContribution    float64
	TotalFees            float64
	TotalOperationalFees float64
	ProductAmounts       map[string]ProductAmounts
}

func (r *repository) GetOJKCustomerReportData(ctx context.Context, filter *crEntity.OJKCustomerReportFilterDto) (*OJKCustomerReportData, error) {
	startDate, endDate, err := parseMonthToDateRange(filter.Month, filter.Year)
	if err != nil {
		return nil, err
	}

	data := &OJKCustomerReportData{}

	customer, err := r.GetCustomerBasicInfo(ctx, filter.CustomerID)
	if err != nil {
		return nil, err
	}
	data.Customer = customer

	summary, err := r.GetInvestmentSummary(ctx, filter.CustomerID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	data.InvestmentItems = summary.Items
	data.ProductAmounts = summary.ProductAmounts
	data.TotalContribution = summary.TotalContribution
	data.TotalFees = summary.TotalFees
	data.TotalOperationalFees = summary.TotalOperationalFees

	navData, err := r.GetNAVData(ctx, endDate)
	if err != nil {
		return nil, err
	}
	data.NAVData = navData

	distributions, err := r.GetInvestmentDistributions(ctx, filter.CustomerID)
	if err != nil {
		return nil, err
	}
	data.InvestmentDistributions = distributions

	openingBalance, err := r.GetOpeningBalance(ctx, filter.CustomerID, startDate)
	if err != nil {
		return nil, err
	}
	data.OpeningBalance = openingBalance

	historicalCost, err := r.GetHistoricalCost(ctx, filter.CustomerID, endDate)
	if err != nil {
		return nil, err
	}
	data.HistoricalCost = historicalCost

	currentPortfolioValue, err := r.GetCurrentPortfolioValue(ctx, filter.CustomerID, endDate)
	if err != nil {
		return nil, err
	}
	data.CurrentPortfolioValue = currentPortfolioValue

	// Get operational fees from transaction_fee table
	operationalFees, err := r.GetOperationalFees(ctx, filter.CustomerID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Add operational fees to product amounts
	for productID, opFee := range operationalFees {
		if amounts, exists := data.ProductAmounts[productID]; exists {
			amounts.OperationalFeeAmount = opFee
			data.ProductAmounts[productID] = amounts
			data.TotalOperationalFees += opFee
		}
	}

	return data, nil
}

func (r *repository) GetCustomerBasicInfo(ctx context.Context, customerID string) (*entity.CustomerDto, error) {
	var customer model.Customer

	if err := r.db.WithContext(ctx).
		Preload("Company").
		Where("id = ? AND deleted_at IS NULL", customerID).
		First(&customer).Error; err != nil {
		return nil, err
	}

	return (&entity.CustomerDto{}).FromModel(&customer), nil
}

func (r *repository) GetInvestmentSummary(ctx context.Context, customerID string, startDate, endDate time.Time) (*InvestmentSummaryData, error) {
	var items []model.InvestmentItem

	query := r.db.WithContext(ctx).
		Preload("InvestmentProduct").
		Preload("Customer").
		Preload("Customer.Company").
		Joins("JOIN investment ON investment.id = investment_item.investment_id").
		Where("investment_item.customer_id = ?", customerID).
		Where("investment_item.investment_at >= ? AND investment_item.investment_at <= ?", startDate, endDate).
		Where("investment.status = ?", model.InvestmentStatusSuccess).
		Where("investment_item.deleted_at IS NULL")

	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}

	itemDtos := make([]*ie.InvestmentItemDto, len(items))
	var totalContribution, totalFees float64
	productAmounts := make(map[string]ProductAmounts)

	for i, item := range items {
		dto := (&ie.InvestmentItemDto{}).FromModel(&item)
		itemDtos[i] = dto

		totalContribution += item.TotalAmount
		totalFees += item.FeeAmount

		if item.InvestmentProduct != nil {
			amounts := productAmounts[item.InvestmentProductID]
			amounts.EmployerAmount += item.EmployerAmount
			amounts.EmployeeAmount += item.EmployeeAmount
			amounts.VoluntaryAmount += item.VoluntaryAmount
			amounts.FeeAmount += item.FeeAmount
			productAmounts[item.InvestmentProductID] = amounts
		}
	}

	return &InvestmentSummaryData{
		Items:             itemDtos,
		TotalContribution: totalContribution,
		TotalFees:         totalFees,
		ProductAmounts:    productAmounts,
	}, nil
}

func (r *repository) GetNAVData(ctx context.Context, endDate time.Time) (map[string]*ie.NetAssetValueDto, error) {
	var navList []model.NetAssetValue

	if err := r.db.WithContext(ctx).
		Preload("InvestmentProduct").
		Where("DATE(created_date) = ?", endDate.Format("2006-01-02")).
		Where("deleted_at IS NULL").
		Find(&navList).Error; err != nil {
		return nil, err
	}

	navByProductID := make(map[string]*ie.NetAssetValueDto)
	for _, nav := range navList {
		navByProductID[nav.InvestmentProductID] = (&ie.NetAssetValueDto{}).FromModel(&nav)
	}

	return navByProductID, nil
}

func (r *repository) GetInvestmentDistributions(ctx context.Context, customerID string) ([]*ie.InvestmentDistributionDto, error) {
	var distributions []model.InvestmentDistribution

	// Get company ID from customer
	var customer model.Customer
	if err := r.db.Where("id = ?", customerID).First(&customer).Error; err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}

	// Get distributions based on customer company status
	var query *gorm.DB
	if customer.CompanyID != "" {
		// Customer with company - get company-level distributions
		query = r.db.WithContext(ctx).
			Preload("InvestmentProduct").
			Where("company_id = ?", customer.CompanyID).
			Where("deleted_at IS NULL")
	} else {
		// Customer without company - get customer-level distributions
		query = r.db.WithContext(ctx).
			Preload("InvestmentProduct").
			Where("customer_id = ?", customerID).
			Where("deleted_at IS NULL")
	}

	if err := query.Find(&distributions).Error; err != nil {
		return nil, err
	}

	distributionDtos := make([]*ie.InvestmentDistributionDto, len(distributions))
	for i, dist := range distributions {
		distributionDtos[i] = (&ie.InvestmentDistributionDto{}).FromModel(&dist)
	}

	return distributionDtos, nil
}

func (r *repository) GetOpeningBalance(ctx context.Context, customerID string, beforeDate time.Time) (map[string]float64, error) {
	var unitLinks []model.UnitLink

	if err := r.db.WithContext(ctx).
		Preload("InvestmentProduct").
		Joins("JOIN participant ON participant.id = unit_link.participant_id").
		Where("participant.customer_id = ?", customerID).
		Where("unit_link.transaction_date < ?", beforeDate).
		Where("unit_link.deleted_at IS NULL").
		Find(&unitLinks).Error; err != nil {
		return nil, err
	}

	openingBalance := make(map[string]float64)
	for _, unitLink := range unitLinks {
		openingBalance[unitLink.InvestmentProductID] += unitLink.Ip
	}

	return openingBalance, nil
}

func (r *repository) GetHistoricalCost(ctx context.Context, customerID string, endDate time.Time) (float64, error) {
	var totalCost float64

	if err := r.db.WithContext(ctx).
		Model(&model.UnitLink{}).
		Select("COALESCE(SUM(total_amount), 0)").
		Joins("JOIN participant ON participant.id = unit_link.participant_id").
		Where("participant.customer_id = ?", customerID).
		Where("unit_link.transaction_date <= ?", endDate).
		Where("unit_link.deleted_at IS NULL").
		Scan(&totalCost).Error; err != nil {
		return 0, err
	}

	return totalCost, nil
}

func (r *repository) GetCurrentPortfolioValue(ctx context.Context, customerID string, endDate time.Time) (float64, error) {
	// Get current portfolio units per product for the customer
	var portfolioData []struct {
		InvestmentProductID string
		TotalUnits          float64
	}

	if err := r.db.WithContext(ctx).
		Model(&model.UnitLink{}).
		Select("investment_product_id, COALESCE(SUM(ip), 0) as total_units").
		Joins("JOIN participant ON participant.id = unit_link.participant_id").
		Where("participant.customer_id = ?", customerID).
		Where("unit_link.transaction_date <= ?", endDate).
		Where("unit_link.deleted_at IS NULL").
		Group("investment_product_id").
		Scan(&portfolioData).Error; err != nil {
		return 0, err
	}

	// Get current NAV for calculation
	navData, err := r.GetNAVData(ctx, endDate)
	if err != nil {
		return 0, err
	}

	var totalValue float64
	for _, portfolio := range portfolioData {
		if nav, exists := navData[portfolio.InvestmentProductID]; exists {
			totalValue += portfolio.TotalUnits * nav.Amount
		}
	}

	return totalValue, nil
}

func parseMonthToDateRange(month, year int) (startDate, endDate time.Time, err error) {
	if month < 1 || month > 12 {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid month: %d. Month must be between 1 and 12", month)
	}
	if year < 2000 || year > 2100 {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid year: %d. Year must be between 2000 and 2100", year)
	}

	startDate = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate = startDate.AddDate(0, 1, -1)

	return startDate, endDate, nil
}

// ProductAmounts holds aggregated amounts for a product
type ProductAmounts struct {
	EmployerAmount       float64
	EmployeeAmount       float64
	VoluntaryAmount      float64
	FeeAmount            float64
	OperationalFeeAmount float64
}

func (r *repository) GetOperationalFees(ctx context.Context, customerID string, startDate, endDate time.Time) (map[string]float64, error) {
	var results []struct {
		InvestmentProductID string
		TotalOperationalFee float64
	}

	err := r.db.WithContext(ctx).
		Table("transaction_fee").
		Select("investment_product_id, COALESCE(SUM(operation_fee), 0) as total_operational_fee").
		Joins("JOIN participant ON participant.id = transaction_fee.participant_id").
		Where("participant.customer_id = ?", customerID).
		Where("transaction_date >= ? AND transaction_date <= ?", startDate, endDate).
		Where("transaction_fee.deleted_at IS NULL").
		Group("investment_product_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	operationalFees := make(map[string]float64)
	for _, result := range results {
		operationalFees[result.InvestmentProductID] = result.TotalOperationalFee
	}

	return operationalFees, nil
}
