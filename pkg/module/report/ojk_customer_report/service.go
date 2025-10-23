package ojkcustomerreport

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	crEntity "github.com/raymondsugiarto/coffee-api/pkg/entity/customer_report"
	ie "github.com/raymondsugiarto/coffee-api/pkg/entity/investment"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/module/customer"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

type OJKCustomerReportService interface {
	GenerateExcelOJKCustomerReportBytes(ctx context.Context, filter *crEntity.OJKCustomerReportFilterDto) ([]byte, error)
	GenerateExcelForCompany(ctx context.Context, filter *crEntity.OJKCustomerReportFilterDto) ([]byte, error)
	GetOJKCustomerReportData(ctx context.Context, filter *crEntity.OJKCustomerReportFilterDto) (*crEntity.OJKCustomerReportDataDto, error)
}

type ojkCustomerReportService struct {
	repository      Repository
	customerService customer.Service
	excelGenerator  *ExcelGenerator
}

func NewOJKCustomerReportService(repository Repository, customerService customer.Service) OJKCustomerReportService {
	service := &ojkCustomerReportService{
		repository:      repository,
		customerService: customerService,
	}
	service.excelGenerator = NewExcelGenerator()

	return service
}

func (s *ojkCustomerReportService) parseMonthToDateRange(month, year int) (startDate, endDate time.Time, err error) {
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

func (s *ojkCustomerReportService) ValidateDateRange(startDate, endDate time.Time) error {
	if endDate.Before(startDate) {
		return fmt.Errorf("end date cannot be before start date")
	}

	yearsDiff := endDate.Year() - startDate.Year()
	monthsDiff := int(endDate.Month()) - int(startDate.Month())
	totalMonthsDiff := yearsDiff*12 + monthsDiff

	if totalMonthsDiff > MaxReportMonths {
		return fmt.Errorf("date range cannot exceed %d months", MaxReportMonths)
	}

	return nil
}

func (s *ojkCustomerReportService) transformToReportData(data *OJKCustomerReportData, startDate, endDate time.Time) (*crEntity.OJKCustomerReportDataDto, error) {
	summary := s.calculateSummary(data.TotalContribution, data.TotalFees, data.HistoricalCost, data.CurrentPortfolioValue)

	// Default for customer without company
	companyType := model.CompanyType("") // Empty string for no company
	if data.Customer.Company != nil {
		companyType = model.CompanyType(data.Customer.Company.CompanyType)
	}

	transactionSections, err := s.groupTransactionsByFundType(data.InvestmentItems, data.ProductAmounts, data.NAVData, data.InvestmentDistributions, data.OpeningBalance, endDate, companyType)
	if err != nil {
		return nil, err
	}

	// Filter out transaction sections with 0% percentage
	filteredTransactionSections := s.filterTransactionSectionsWithZeroPercentage(transactionSections)

	contributionSummary := s.calculateContributionSummary(data.InvestmentItems)
	feeSummary := s.calculateFeeSummary(data.InvestmentItems, data.TotalOperationalFees)

	// Calculate additional fields
	reportTitle := s.calculateReportTitle(data.Customer)
	totalContribution := s.calculateTotalContribution(contributionSummary, companyType)
	filteredTransactionSections = s.addLastUnitPriceDates(filteredTransactionSections)

	return &crEntity.OJKCustomerReportDataDto{
		Customer: data.Customer,
		Period: &crEntity.ReportPeriodDto{
			StartDate: startDate,
			EndDate:   endDate,
		},
		Summary:             summary,
		TransactionSections: filteredTransactionSections,
		ContributionSummary: contributionSummary,
		FeeSummary:          feeSummary,
		ReportTitle:         reportTitle,
		TotalContribution:   totalContribution,
	}, nil
}

func (s *ojkCustomerReportService) groupTransactionsByFundType(items []*ie.InvestmentItemDto, productAmounts map[string]ProductAmounts, navByProductID map[string]*ie.NetAssetValueDto, distributions []*ie.InvestmentDistributionDto, openingBalance map[string]float64, endDate time.Time, companyType model.CompanyType) ([]*crEntity.TransactionSectionDto, error) {
	if err := s.validateNAVData(navByProductID, endDate); err != nil {
		return nil, err
	}

	if err := s.validateDistributions(distributions, navByProductID, endDate); err != nil {
		return nil, err
	}

	productMap := s.extractUniqueProducts(items)
	distributionByProductID := s.mapDistributionsByProductID(distributions)

	return s.buildTransactionSections(productMap, productAmounts, distributionByProductID, navByProductID, openingBalance, companyType, endDate)
}

func (s *ojkCustomerReportService) validateNAVData(navByProductID map[string]*ie.NetAssetValueDto, endDate time.Time) error {
	if len(navByProductID) == 0 {
		return fmt.Errorf("data NAV tidak tersedia untuk bulan %s. Silakan input data NAV untuk tanggal %s",
			endDate.Format(DateFormatMonthYear), endDate.Format(DateFormatDisplay))
	}
	return nil
}

func (s *ojkCustomerReportService) validateDistributions(distributions []*ie.InvestmentDistributionDto, navByProductID map[string]*ie.NetAssetValueDto, endDate time.Time) error {
	var missingNavProducts []string
	for _, dist := range distributions {
		if _, exists := navByProductID[dist.InvestmentProductID]; !exists {
			productName := "Unknown Product"
			if dist.InvestmentProduct != nil {
				productName = dist.InvestmentProduct.Name
			}
			missingNavProducts = append(missingNavProducts, productName)
		}
	}

	if len(missingNavProducts) > 0 {
		return fmt.Errorf("data NAV tidak tersedia untuk produk investasi berikut pada tanggal %s: %v. Silakan input data NAV untuk semua produk yang terdaftar dalam investment distribution",
			endDate.Format(DateFormatDisplay), missingNavProducts)
	}

	return nil
}

func (s *ojkCustomerReportService) extractUniqueProducts(items []*ie.InvestmentItemDto) map[string]*ie.InvestmentProductDto {
	productMap := make(map[string]*ie.InvestmentProductDto)
	for _, item := range items {
		if item.InvestmentProduct != nil {
			productMap[item.InvestmentProductID] = item.InvestmentProduct
		}
	}
	return productMap
}

func (s *ojkCustomerReportService) mapDistributionsByProductID(distributions []*ie.InvestmentDistributionDto) map[string]*ie.InvestmentDistributionDto {
	distributionByProductID := make(map[string]*ie.InvestmentDistributionDto)
	for _, dist := range distributions {
		distributionByProductID[dist.InvestmentProductID] = dist
	}
	return distributionByProductID
}

func (s *ojkCustomerReportService) buildTransactionSections(
	productMap map[string]*ie.InvestmentProductDto,
	productAmounts map[string]ProductAmounts,
	distributionByProductID map[string]*ie.InvestmentDistributionDto,
	navByProductID map[string]*ie.NetAssetValueDto,
	openingBalance map[string]float64,
	companyType model.CompanyType,
	endDate time.Time,
) ([]*crEntity.TransactionSectionDto, error) {
	var sections []*crEntity.TransactionSectionDto
	fundTypeMap := make(map[string][]*crEntity.TransactionItemDto)
	fundBalanceMap := make(map[string]float64)

	for productID, product := range productMap {
		fundType := product.Name
		nav, exists := navByProductID[productID]
		if !exists {
			return nil, fmt.Errorf("NAV tidak tersedia untuk produk %s pada tanggal %s",
				product.Name, endDate.Format(DateFormatDisplay))
		}

		transactions := s.buildProductTransactions(product, nav, productAmounts[productID], openingBalance[productID], &fundBalanceMap, companyType, endDate)
		fundTypeMap[fundType] = transactions
	}

	for fundType, transactions := range fundTypeMap {
		productID := s.findProductIDByName(productMap, fundType)
		percentage := s.getDistributionPercentage(distributionByProductID, productID)
		finalBalance := fundBalanceMap[fundType]
		finalUnitPrice := navByProductID[productID].Amount
		finalValue := finalBalance * finalUnitPrice

		sections = append(sections, &crEntity.TransactionSectionDto{
			FundType:     fundType,
			Percentage:   percentage,
			Transactions: transactions,
			FinalBalance: finalBalance,
			UnitPrice:    finalUnitPrice,
			FinalValue:   finalValue,
		})
	}

	return sections, nil
}

func (s *ojkCustomerReportService) buildProductTransactions(
	product *ie.InvestmentProductDto,
	nav *ie.NetAssetValueDto,
	amounts ProductAmounts,
	openingBalanceUnits float64,
	fundBalanceMap *map[string]float64,
	companyType model.CompanyType,
	endDate time.Time,
) []*crEntity.TransactionItemDto {
	fundType := product.Name
	unitPrice := nav.Amount
	transactionDate := endDate

	// Initialize with opening balance from database
	openingBalanceAmount := openingBalanceUnits * unitPrice
	(*fundBalanceMap)[fundType] = openingBalanceUnits

	transactions := []*crEntity.TransactionItemDto{
		{
			TransactionType:  TransactionTypeOpeningBalance,
			TransactionDate:  transactionDate,
			InvestmentValue:  openingBalanceAmount,
			UnitPriceDate:    transactionDate,
			UnitPrice:        unitPrice,
			TransactionUnits: openingBalanceUnits,
			BalanceUnits:     (*fundBalanceMap)[fundType],
		},
	}

	transactionMap := s.buildTransactionMap(amounts, companyType)

	for _, txType := range TransactionOrder {
		amount, exists := transactionMap[txType]
		if !exists {
			continue
		}

		isFee := txType == TransactionTypeAdminFee || txType == TransactionTypeOperationalFee
		transactionUnits := s.calculateTransactionUnits(amount, unitPrice, isFee)
		(*fundBalanceMap)[fundType] += transactionUnits

		transactions = append(transactions, &crEntity.TransactionItemDto{
			TransactionType:  txType,
			TransactionDate:  transactionDate,
			InvestmentValue:  amount,
			UnitPriceDate:    transactionDate,
			UnitPrice:        unitPrice,
			TransactionUnits: transactionUnits,
			BalanceUnits:     (*fundBalanceMap)[fundType],
		})
	}

	return transactions
}

func (s *ojkCustomerReportService) findProductIDByName(productMap map[string]*ie.InvestmentProductDto, fundType string) string {
	for pid, product := range productMap {
		if product.Name == fundType {
			return pid
		}
	}
	return ""
}

func (s *ojkCustomerReportService) getDistributionPercentage(distributionByProductID map[string]*ie.InvestmentDistributionDto, productID string) int {
	if dist, exists := distributionByProductID[productID]; exists {
		return int(dist.Percent)
	}
	return 0
}

func (s *ojkCustomerReportService) GenerateExcelOJKCustomerReportBytes(ctx context.Context, filter *crEntity.OJKCustomerReportFilterDto) ([]byte, error) {
	startDate, endDate, err := s.parseMonthToDateRange(filter.Month, filter.Year)
	if err != nil {
		return nil, err
	}

	if err := s.ValidateDateRange(startDate, endDate); err != nil {
		return nil, err
	}

	data, err := s.repository.GetOJKCustomerReportData(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get OJK customer report data: %w", err)
	}

	reportData, err := s.transformToReportData(data, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return s.excelGenerator.GenerateExcel(reportData)
}

func (s *ojkCustomerReportService) filterTransactionSectionsWithZeroPercentage(sections []*crEntity.TransactionSectionDto) []*crEntity.TransactionSectionDto {
	var filtered []*crEntity.TransactionSectionDto
	for _, section := range sections {
		if section.Percentage > 0 {
			filtered = append(filtered, section)
		}
	}
	return filtered
}

// Calculation methods
func (s *ojkCustomerReportService) calculateSummary(totalContribution, totalFees, historicalCost, currentPortfolioValue float64) *crEntity.OJKCustomerReportSummaryDto {
	developmentResults := currentPortfolioValue - historicalCost
	managedFundValue := currentPortfolioValue

	return &crEntity.OJKCustomerReportSummaryDto{
		AccumulatedContribution:       totalContribution,
		AccumulatedDevelopmentResults: developmentResults,
		AccumulatedFees:               totalFees,
		ManagedFundValue:              managedFundValue,
	}
}

func (s *ojkCustomerReportService) calculateContributionSummary(items []*ie.InvestmentItemDto) *crEntity.ContributionSummaryDto {
	summary := &crEntity.ContributionSummaryDto{}

	for _, item := range items {
		summary.EmployerContribution += item.EmployerAmount
		summary.EmployeeContribution += item.EmployeeAmount
		summary.VoluntaryContribution += item.VoluntaryAmount
		summary.EducationFund += item.EducationFundAmount
	}

	return summary
}

func (s *ojkCustomerReportService) calculateFeeSummary(items []*ie.InvestmentItemDto, totalOperationalFees float64) *crEntity.FeeSummaryDto {
	var adminFee float64

	for _, item := range items {
		adminFee += item.FeeAmount
	}

	return &crEntity.FeeSummaryDto{
		AdministrationFee: adminFee,
		OperationalFee:    totalOperationalFees,
	}
}

func (s *ojkCustomerReportService) calculateTransactionUnits(amount, unitPrice float64, isFee bool) float64 {
	if unitPrice <= 0 {
		return 0
	}

	units := amount / unitPrice
	if isFee {
		units = -units
	}

	return units
}

func (s *ojkCustomerReportService) buildTransactionMap(amounts ProductAmounts, companyType model.CompanyType) map[string]float64 {
	transactionMap := map[string]float64{}

	// Skip employer & voluntary contributions for customer without company
	if companyType != "" {
		// Customer with company
		transactionMap[TransactionTypeEmployerContrib] = amounts.EmployerAmount
		if companyType == model.CompanyTypePPIP {
			transactionMap[TransactionTypeVoluntaryContrib] = amounts.VoluntaryAmount
		}
	}

	// Employee contribution is always present for all customers
	transactionMap[TransactionTypeEmployeeContrib] = amounts.EmployeeAmount

	transactionMap[TransactionTypeAdminFee] = amounts.FeeAmount
	transactionMap[TransactionTypeOperationalFee] = amounts.OperationalFeeAmount

	return transactionMap
}

func (s *ojkCustomerReportService) GenerateExcelForCompany(ctx context.Context, filter *crEntity.OJKCustomerReportFilterDto) ([]byte, error) {
	// Get company ID from context
	companyID := shared.GetCompanyID(ctx)
	if companyID == nil {
		return nil, status.New(status.Unauthorized, errors.New("company ID not found in context"))
	}

	// Validate that the customer belongs to the requesting company
	if filter.CustomerID != "" {
		// Get customer info to validate company ownership
		customerInfo, err := s.customerService.FindByID(ctx, filter.CustomerID)
		if err != nil {
			return nil, status.New(status.BadRequest, errors.New("customer not found"))
		}

		// Check if customer belongs to the requesting company
		if customerInfo.CompanyID == "" || customerInfo.CompanyID != *companyID {
			return nil, status.New(status.Forbidden, errors.New("customer does not belong to your company"))
		}
	}

	// Call the existing method to generate the report
	return s.GenerateExcelOJKCustomerReportBytes(ctx, filter)
}

func (s *ojkCustomerReportService) GetOJKCustomerReportData(ctx context.Context, filter *crEntity.OJKCustomerReportFilterDto) (*crEntity.OJKCustomerReportDataDto, error) {
	startDate, endDate, err := s.parseMonthToDateRange(filter.Month, filter.Year)
	if err != nil {
		return nil, err
	}

	if err := s.ValidateDateRange(startDate, endDate); err != nil {
		return nil, err
	}

	data, err := s.repository.GetOJKCustomerReportData(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get OJK customer report data: %w", err)
	}

	reportData, err := s.transformToReportData(data, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to transform report data: %w", err)
	}

	return reportData, nil
}

func (s *ojkCustomerReportService) calculateReportTitle(customer *entity.CustomerDto) string {
	if customer.Company != nil {
		return "LAPORAN TRANSAKSI PESERTA - PPIP KUMPULAN"
	}
	return "LAPORAN TRANSAKSI PESERTA - PPIP MANDIRI"
}

func (s *ojkCustomerReportService) calculateTotalContribution(summary *crEntity.ContributionSummaryDto, companyType model.CompanyType) float64 {
	if companyType == "" {
		// Individual customer - only employee contribution
		return summary.EmployeeContribution
	}
	if companyType == model.CompanyTypeDKP {
		return summary.EmployerContribution
	}
	// PPIP customer with company
	return summary.EmployerContribution + summary.EmployeeContribution + summary.VoluntaryContribution + summary.EducationFund
}

// addLastUnitPriceDates adds last unit price date to each transaction section
func (s *ojkCustomerReportService) addLastUnitPriceDates(sections []*crEntity.TransactionSectionDto) []*crEntity.TransactionSectionDto {
	for _, section := range sections {
		if len(section.Transactions) > 0 {
			section.LastUnitPriceDate = section.Transactions[len(section.Transactions)-1].UnitPriceDate
		}
	}
	return sections
}
