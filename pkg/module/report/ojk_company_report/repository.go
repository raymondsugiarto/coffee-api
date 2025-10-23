package ojkcompanyreport

import (
	"context"
	"fmt"
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	crEntity "github.com/raymondsugiarto/coffee-api/pkg/entity/company_report"
	ie "github.com/raymondsugiarto/coffee-api/pkg/entity/investment"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"gorm.io/gorm"
)

type Repository interface {
	GetOJKCompanyReportData(ctx context.Context, filter *crEntity.OJKCompanyReportFilterDto) (*OJKCompanyReportData, error)
	GetCompanyBasicInfo(ctx context.Context, companyID string) (*entity.CompanyDto, error)
	GetInvestmentSummary(ctx context.Context, companyID string, startDate, endDate time.Time) (*InvestmentSummaryData, error)
	GetParticipantCount(ctx context.Context, companyID string) (int, error)
	GetNAVData(ctx context.Context, endDate time.Time) (map[string]*ie.NetAssetValueDto, error)
	GetInvestmentDistributions(ctx context.Context, companyID string) ([]*ie.InvestmentDistributionDto, error)
	GetOpeningBalance(ctx context.Context, companyID string, beforeDate time.Time) (map[string]float64, error)
	GetHistoricalCost(ctx context.Context, companyID string, endDate time.Time) (float64, error)
	GetCurrentPortfolioValue(ctx context.Context, companyID string, endDate time.Time) (float64, error)
	GetOperationalFees(ctx context.Context, companyID string, startDate, endDate time.Time) (map[string]float64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

type OJKCompanyReportData struct {
	Company                 *entity.CompanyDto
	InvestmentItems         []*ie.InvestmentItemDto
	ParticipantCount        int
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

func (r *repository) GetOJKCompanyReportData(ctx context.Context, filter *crEntity.OJKCompanyReportFilterDto) (*OJKCompanyReportData, error) {
	startDate, endDate, err := parseMonthToDateRange(filter.Month, filter.Year)
	if err != nil {
		return nil, err
	}

	data := &OJKCompanyReportData{}

	company, err := r.GetCompanyBasicInfo(ctx, filter.CompanyID)
	if err != nil {
		return nil, err
	}
	data.Company = company

	summary, err := r.GetInvestmentSummary(ctx, filter.CompanyID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	data.InvestmentItems = summary.Items
	data.ProductAmounts = summary.ProductAmounts
	data.TotalContribution = summary.TotalContribution
	data.TotalFees = summary.TotalFees
	data.TotalOperationalFees = summary.TotalOperationalFees

	participantCount, err := r.GetParticipantCount(ctx, filter.CompanyID)
	if err != nil {
		return nil, err
	}
	data.ParticipantCount = participantCount

	navData, err := r.GetNAVData(ctx, endDate)
	if err != nil {
		return nil, err
	}
	data.NAVData = navData

	distributions, err := r.GetInvestmentDistributions(ctx, filter.CompanyID)
	if err != nil {
		return nil, err
	}
	data.InvestmentDistributions = distributions

	openingBalance, err := r.GetOpeningBalance(ctx, filter.CompanyID, startDate)
	if err != nil {
		return nil, err
	}
	data.OpeningBalance = openingBalance

	historicalCost, err := r.GetHistoricalCost(ctx, filter.CompanyID, endDate)
	if err != nil {
		return nil, err
	}
	data.HistoricalCost = historicalCost

	currentPortfolioValue, err := r.GetCurrentPortfolioValue(ctx, filter.CompanyID, endDate)
	if err != nil {
		return nil, err
	}
	data.CurrentPortfolioValue = currentPortfolioValue

	// Get operational fees from transaction_fee table
	operationalFees, err := r.GetOperationalFees(ctx, filter.CompanyID, startDate, endDate)
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

func (r *repository) GetCompanyBasicInfo(ctx context.Context, companyID string) (*entity.CompanyDto, error) {
	var company model.Company

	if err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", companyID).
		First(&company).Error; err != nil {
		return nil, err
	}

	return (&entity.CompanyDto{}).FromModel(&company), nil
}

func (r *repository) GetInvestmentSummary(ctx context.Context, companyID string, startDate, endDate time.Time) (*InvestmentSummaryData, error) {
	var items []model.InvestmentItem

	query := r.db.WithContext(ctx).
		Preload("InvestmentProduct").
		Preload("Customer").
		Preload("Customer.Company").
		Joins("JOIN customer ON customer.id = investment_item.customer_id").
		Joins("JOIN investment ON investment.id = investment_item.investment_id").
		Where("customer.company_id = ?", companyID).
		Where("investment_item.investment_at >= ? AND investment_item.investment_at <= ?", startDate, endDate).
		// Where("investment_item.status = ?", model.InvestmentStatusSuccess).  NOTE: investment items status not updated to success
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

func (r *repository) GetParticipantCount(ctx context.Context, companyID string) (int, error) {
	var count int64

	if err := r.db.WithContext(ctx).
		Model(&model.Participant{}).
		Joins("JOIN customer ON customer.id = participant.customer_id").
		Where("customer.company_id = ?", companyID).
		Where("participant.status = ?", model.ParticipantStatusActive).
		Where("participant.deleted_at IS NULL").
		Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
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

func (r *repository) GetInvestmentDistributions(ctx context.Context, companyID string) ([]*ie.InvestmentDistributionDto, error) {
	var distributions []model.InvestmentDistribution

	if err := r.db.WithContext(ctx).
		Preload("InvestmentProduct").
		Where("company_id = ?", companyID).
		Where("deleted_at IS NULL").
		Find(&distributions).Error; err != nil {
		return nil, err
	}

	distributionDtos := make([]*ie.InvestmentDistributionDto, len(distributions))
	for i, dist := range distributions {
		distributionDtos[i] = (&ie.InvestmentDistributionDto{}).FromModel(&dist)
	}

	return distributionDtos, nil
}

func (r *repository) GetOpeningBalance(ctx context.Context, companyID string, beforeDate time.Time) (map[string]float64, error) {
	var unitLinks []model.UnitLink

	if err := r.db.WithContext(ctx).
		Preload("InvestmentProduct").
		Joins("JOIN participant ON participant.id = unit_link.participant_id").
		Joins("JOIN customer ON customer.id = participant.customer_id").
		Where("customer.company_id = ?", companyID).
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

func (r *repository) GetHistoricalCost(ctx context.Context, companyID string, endDate time.Time) (float64, error) {
	var totalCost float64

	if err := r.db.WithContext(ctx).
		Model(&model.UnitLink{}).
		Select("COALESCE(SUM(total_amount), 0)").
		Joins("JOIN participant ON participant.id = unit_link.participant_id").
		Joins("JOIN customer ON customer.id = participant.customer_id").
		Where("customer.company_id = ?", companyID).
		Where("unit_link.transaction_date <= ?", endDate).
		Where("unit_link.deleted_at IS NULL").
		Scan(&totalCost).Error; err != nil {
		return 0, err
	}

	return totalCost, nil
}

func (r *repository) GetCurrentPortfolioValue(ctx context.Context, companyID string, endDate time.Time) (float64, error) {
	// Get current portfolio units per product
	var portfolioData []struct {
		InvestmentProductID string
		TotalUnits          float64
	}

	if err := r.db.WithContext(ctx).
		Model(&model.UnitLink{}).
		Select("investment_product_id, COALESCE(SUM(ip), 0) as total_units").
		Joins("JOIN participant ON participant.id = unit_link.participant_id").
		Joins("JOIN customer ON customer.id = participant.customer_id").
		Where("customer.company_id = ?", companyID).
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

func (r *repository) GetOperationalFees(ctx context.Context, companyID string, startDate, endDate time.Time) (map[string]float64, error) {
	var results []struct {
		InvestmentProductID string
		TotalOperationalFee float64
	}

	err := r.db.WithContext(ctx).
		Table("transaction_fee").
		Select("investment_product_id, COALESCE(SUM(operation_fee), 0) as total_operational_fee").
		Where("company_id = ?", companyID).
		Where("transaction_date >= ? AND transaction_date <= ?", startDate, endDate).
		Where("deleted_at IS NULL").
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
