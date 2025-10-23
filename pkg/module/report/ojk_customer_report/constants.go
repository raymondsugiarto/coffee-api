package ojkcustomerreport

const (
	// Date and time formats
	DateFormatMonthYear = "2006-01"
	DateFormatDisplay   = "02/01/2006"

	// Validation constraints
	MaxReportMonths = 12

	// Transaction types
	TransactionTypeOpeningBalance   = "Saldo Awal"
	TransactionTypeEmployerContrib  = "Iuran Pemberi Kerja"
	TransactionTypeEmployeeContrib  = "Iuran Pekerja"
	TransactionTypeVoluntaryContrib = "Iuran Sukarela Pekerja"
	TransactionTypeAdminFee         = "Biaya Admin Iuran"
	TransactionTypeOperationalFee   = "Biaya Pengelolaan Investasi"

	// Distribution type
	DistributionTypeCustomer = "customer"
)

// TransactionOrder defines the order of transactions in reports
var TransactionOrder = []string{
	TransactionTypeEmployerContrib,
	TransactionTypeEmployeeContrib,
	TransactionTypeVoluntaryContrib,
	TransactionTypeAdminFee,
	TransactionTypeOperationalFee,
}
