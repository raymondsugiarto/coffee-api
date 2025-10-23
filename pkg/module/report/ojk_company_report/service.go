package ojkcompanyreport

import (
	"context"
	"fmt"
	"time"

	crEntity "github.com/raymondsugiarto/coffee-api/pkg/entity/company_report"
	ie "github.com/raymondsugiarto/coffee-api/pkg/entity/investment"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
)

type OJKCompanyReportService interface {
	GenerateExcelOJKCompanyReportBytes(ctx context.Context, filter *crEntity.OJKCompanyReportFilterDto) ([]byte, error)
	GetOJKCompanyReportData(ctx context.Context, filter *crEntity.OJKCompanyReportFilterDto) (*crEntity.OJKCompanyReportDataDto, error)
}

type ojkCompanyReportService struct {
	repository     Repository
	excelGenerator *ExcelGenerator
}

func NewOJKCompanyReportService(repository Repository) OJKCompanyReportService {
	service := &ojkCompanyReportService{
		repository: repository,
	}
	service.excelGenerator = NewExcelGenerator()

	return service
}

func (s *ojkCompanyReportService) parseMonthToDateRange(month, year int) (startDate, endDate time.Time, err error) {
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

func (s *ojkCompanyReportService) ValidateDateRange(startDate, endDate time.Time) error {
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

func (s *ojkCompanyReportService) transformToReportData(data *OJKCompanyReportData, startDate, endDate time.Time) (*crEntity.OJKCompanyReportDataDto, error) {
	summary := s.calculateSummary(data.ParticipantCount, data.TotalContribution, data.TotalFees+data.TotalOperationalFees, data.HistoricalCost, data.CurrentPortfolioValue)

	transactionSections, err := s.groupTransactionsByFundType(data.InvestmentItems, data.ProductAmounts, data.NAVData, data.InvestmentDistributions, data.OpeningBalance, endDate, model.CompanyType(data.Company.CompanyType))
	if err != nil {
		return nil, err
	}

	// Filter out transaction sections with 0% percentage
	filteredTransactionSections := s.filterTransactionSectionsWithZeroPercentage(transactionSections)

	contributionSummary := s.calculateContributionSummary(data.InvestmentItems)
	feeSummary := s.calculateFeeSummary(data.InvestmentItems, data.TotalOperationalFees)

	// Calculate additional fields
	reportTitle := s.calculateReportTitle(data.Company.CompanyType)
	totalContribution := s.calculateTotalContribution(contributionSummary, data.Company.CompanyType)
	filteredTransactionSections = s.addLastUnitPriceDates(filteredTransactionSections)

	return &crEntity.OJKCompanyReportDataDto{
		Company: data.Company,
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

func (s *ojkCompanyReportService) groupTransactionsByFundType(items []*ie.InvestmentItemDto, productAmounts map[string]ProductAmounts, navByProductID map[string]*ie.NetAssetValueDto, distributions []*ie.InvestmentDistributionDto, openingBalance map[string]float64, endDate time.Time, companyType model.CompanyType) ([]*crEntity.TransactionSectionDto, error) {
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

func (s *ojkCompanyReportService) validateNAVData(navByProductID map[string]*ie.NetAssetValueDto, endDate time.Time) error {
	if len(navByProductID) == 0 {
		return fmt.Errorf("data NAV tidak tersedia untuk bulan %s. Silakan input data NAV untuk tanggal %s",
			endDate.Format(DateFormatMonthYear), endDate.Format(DateFormatDisplay))
	}
	return nil
}

func (s *ojkCompanyReportService) validateDistributions(distributions []*ie.InvestmentDistributionDto, navByProductID map[string]*ie.NetAssetValueDto, endDate time.Time) error {
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

func (s *ojkCompanyReportService) extractUniqueProducts(items []*ie.InvestmentItemDto) map[string]*ie.InvestmentProductDto {
	productMap := make(map[string]*ie.InvestmentProductDto)
	for _, item := range items {
		if item.InvestmentProduct != nil {
			productMap[item.InvestmentProductID] = item.InvestmentProduct
		}
	}
	return productMap
}

func (s *ojkCompanyReportService) mapDistributionsByProductID(distributions []*ie.InvestmentDistributionDto) map[string]*ie.InvestmentDistributionDto {
	distributionByProductID := make(map[string]*ie.InvestmentDistributionDto)
	for _, dist := range distributions {
		distributionByProductID[dist.InvestmentProductID] = dist
	}
	return distributionByProductID
}

func (s *ojkCompanyReportService) buildTransactionSections(
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

func (s *ojkCompanyReportService) buildProductTransactions(
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

		// Make investment value negative for fees
		investmentValue := amount
		if isFee {
			investmentValue = -amount
		}

		transactions = append(transactions, &crEntity.TransactionItemDto{
			TransactionType:  txType,
			TransactionDate:  transactionDate,
			InvestmentValue:  investmentValue,
			UnitPriceDate:    transactionDate,
			UnitPrice:        unitPrice,
			TransactionUnits: transactionUnits,
			BalanceUnits:     (*fundBalanceMap)[fundType],
		})
	}

	return transactions
}

func (s *ojkCompanyReportService) findProductIDByName(productMap map[string]*ie.InvestmentProductDto, fundType string) string {
	for pid, product := range productMap {
		if product.Name == fundType {
			return pid
		}
	}
	return ""
}

func (s *ojkCompanyReportService) getDistributionPercentage(distributionByProductID map[string]*ie.InvestmentDistributionDto, productID string) int {
	if dist, exists := distributionByProductID[productID]; exists {
		return int(dist.Percent)
	}
	return 0
}

func (s *ojkCompanyReportService) GenerateExcelOJKCompanyReportBytes(ctx context.Context, filter *crEntity.OJKCompanyReportFilterDto) ([]byte, error) {
	startDate, endDate, err := s.parseMonthToDateRange(filter.Month, filter.Year)
	if err != nil {
		return nil, err
	}

	if err := s.ValidateDateRange(startDate, endDate); err != nil {
		return nil, err
	}

	data, err := s.repository.GetOJKCompanyReportData(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get OJK company report data: %w", err)
	}

	reportData, err := s.transformToReportData(data, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return s.excelGenerator.GenerateExcel(reportData)
}

func (s *ojkCompanyReportService) filterTransactionSectionsWithZeroPercentage(sections []*crEntity.TransactionSectionDto) []*crEntity.TransactionSectionDto {
	var filtered []*crEntity.TransactionSectionDto
	for _, section := range sections {
		if section.Percentage > 0 {
			filtered = append(filtered, section)
		}
	}
	return filtered
}

// ProductAmounts holds aggregated amounts for a product
type ProductAmounts struct {
	EmployerAmount       float64
	EmployeeAmount       float64
	VoluntaryAmount      float64
	FeeAmount            float64
	OperationalFeeAmount float64
}

// Calculation methods
func (s *ojkCompanyReportService) calculateSummary(participantCount int, totalContribution, totalFees, historicalCost, currentPortfolioValue float64) *crEntity.OJKCompanyReportSummaryDto {
	developmentResults := currentPortfolioValue - historicalCost
	managedFundValue := currentPortfolioValue

	return &crEntity.OJKCompanyReportSummaryDto{
		ParticipantCount:              participantCount,
		AccumulatedContribution:       totalContribution,
		AccumulatedDevelopmentResults: developmentResults,
		AccumulatedFees:               totalFees,
		ManagedFundValue:              managedFundValue,
	}
}

func (s *ojkCompanyReportService) calculateContributionSummary(items []*ie.InvestmentItemDto) *crEntity.ContributionSummaryDto {
	summary := &crEntity.ContributionSummaryDto{}

	for _, item := range items {
		summary.EmployerContribution += item.EmployerAmount
		summary.EmployeeContribution += item.EmployeeAmount
		summary.VoluntaryContribution += item.VoluntaryAmount
		summary.EducationFund += item.EducationFundAmount
	}

	return summary
}

func (s *ojkCompanyReportService) calculateFeeSummary(items []*ie.InvestmentItemDto, totalOperationalFees float64) *crEntity.FeeSummaryDto {
	var adminFee float64

	for _, item := range items {
		adminFee += item.FeeAmount
	}

	return &crEntity.FeeSummaryDto{
		AdministrationFee: adminFee,
		OperationalFee:    totalOperationalFees,
	}
}

func (s *ojkCompanyReportService) calculateTransactionUnits(amount, unitPrice float64, isFee bool) float64 {
	if unitPrice <= 0 {
		return 0
	}

	units := amount / unitPrice
	if isFee {
		units = -units
	}

	return units
}

func (s *ojkCompanyReportService) buildTransactionMap(amounts ProductAmounts, companyType model.CompanyType) map[string]float64 {
	transactionMap := map[string]float64{
		TransactionTypeEmployerContrib: amounts.EmployerAmount,
	}

	if companyType == model.CompanyTypePPIP {
		transactionMap[TransactionTypeEmployeeContrib] = amounts.EmployeeAmount
		transactionMap[TransactionTypeVoluntaryContrib] = amounts.VoluntaryAmount
	}

	transactionMap[TransactionTypeAdminFee] = amounts.FeeAmount
	transactionMap[TransactionTypeOperationalFee] = amounts.OperationalFeeAmount

	return transactionMap
}

func (s *ojkCompanyReportService) GetOJKCompanyReportData(ctx context.Context, filter *crEntity.OJKCompanyReportFilterDto) (*crEntity.OJKCompanyReportDataDto, error) {
	startDate, endDate, err := s.parseMonthToDateRange(filter.Month, filter.Year)
	if err != nil {
		return nil, err
	}

	if err := s.ValidateDateRange(startDate, endDate); err != nil {
		return nil, err
	}

	data, err := s.repository.GetOJKCompanyReportData(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get OJK company report data: %w", err)
	}

	reportData, err := s.transformToReportData(data, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to transform report data: %w", err)
	}

	return reportData, nil
}

func (s *ojkCompanyReportService) calculateReportTitle(companyType string) string {
	switch companyType {
	case "DKP":
		return "LAPORAN TRANSAKSI PERUSAHAAN - DKP"
	case "PPIP":
		return "LAPORAN TRANSAKSI PERUSAHAAN - PPIP KUMPULAN"
	default:
		return "LAPORAN TRANSAKSI PERUSAHAAN - PPIP KUMPULAN"
	}
}

func (s *ojkCompanyReportService) calculateTotalContribution(summary *crEntity.ContributionSummaryDto, companyType string) float64 {
	if companyType == "DKP" {
		return summary.EmployerContribution
	}
	return summary.EmployerContribution + summary.EmployeeContribution + summary.VoluntaryContribution + summary.EducationFund
}

func (s *ojkCompanyReportService) addLastUnitPriceDates(sections []*crEntity.TransactionSectionDto) []*crEntity.TransactionSectionDto {
	for _, section := range sections {
		if len(section.Transactions) > 0 {
			section.LastUnitPriceDate = section.Transactions[len(section.Transactions)-1].UnitPriceDate
		}
	}
	return sections
}
