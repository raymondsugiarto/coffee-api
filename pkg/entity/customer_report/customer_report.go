package customer_report

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
)

// OJKCustomerReportFilterDto defines the query filter for generating an OJK customer report.
type OJKCustomerReportFilterDto struct {
	CustomerID string `json:"customerId" validate:"required"`
	Month      int    `json:"month" validate:"required,min=1,max=12"`
	Year       int    `json:"year" validate:"required,min=2000,max=2100"`
}

// OJKCustomerReportResponseDto defines the response after generating an OJK customer report.
type OJKCustomerReportResponseDto struct {
	FilePath string `json:"filePath"` // For both PDF and Excel files
}

// OJKCustomerReportDataDto contains all the data needed to generate the OJK customer report Excel.
type OJKCustomerReportDataDto struct {
	Customer            *entity.CustomerDto          `json:"customer"`
	Period              *ReportPeriodDto             `json:"period"`
	Summary             *OJKCustomerReportSummaryDto `json:"summary"`
	TransactionSections []*TransactionSectionDto     `json:"transactionSections"`
	ContributionSummary *ContributionSummaryDto      `json:"contributionSummary"`
	FeeSummary          *FeeSummaryDto               `json:"feeSummary"`
	ReportTitle         string                       `json:"reportTitle"`
	TotalContribution   float64                      `json:"totalContribution"`
}

// ReportPeriodDto represents the reporting period.
type ReportPeriodDto struct {
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}

// OJKCustomerReportSummaryDto represents the overall summary section for individual customer.
type OJKCustomerReportSummaryDto struct {
	AccumulatedContribution       float64 `json:"accumulatedContribution"`
	AccumulatedDevelopmentResults float64 `json:"accumulatedDevelopmentResults"`
	AccumulatedFees               float64 `json:"accumulatedFees"`
	ManagedFundValue              float64 `json:"managedFundValue"`
}

// TransactionSectionDto represents a section for each investment fund type.
type TransactionSectionDto struct {
	FundType          string                `json:"fundType"`
	Percentage        int                   `json:"percentage"`
	Transactions      []*TransactionItemDto `json:"transactions"`
	FinalBalance      float64               `json:"finalBalance"`
	UnitPrice         float64               `json:"unitPrice"`
	FinalValue        float64               `json:"finalValue"`
	LastUnitPriceDate time.Time             `json:"lastUnitPriceDate"`
}

// TransactionItemDto represents individual transaction entries.
type TransactionItemDto struct {
	TransactionType  string    `json:"transactionType"`
	TransactionDate  time.Time `json:"transactionDate"`
	InvestmentValue  float64   `json:"investmentValue"`
	UnitPriceDate    time.Time `json:"unitPriceDate"`
	UnitPrice        float64   `json:"unitPrice"`
	TransactionUnits float64   `json:"transactionUnits"`
	BalanceUnits     float64   `json:"balanceUnits"`
}

// ContributionSummaryDto represents the contribution summary section.
type ContributionSummaryDto struct {
	EmployerContribution  float64 `json:"employerContribution"`
	EmployeeContribution  float64 `json:"employeeContribution"`
	VoluntaryContribution float64 `json:"voluntaryContribution"`
	EducationFund         float64 `json:"educationFund"`
}

// FeeSummaryDto represents the fee summary section.
type FeeSummaryDto struct {
	AdministrationFee float64 `json:"administrationFee"`
	OperationalFee    float64 `json:"operationalFee"`
}
