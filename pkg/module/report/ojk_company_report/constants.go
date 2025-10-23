package ojkcompanyreport

const (
	// Date and time formats
	DateFormatMonthYear = "2006-01"
	DateFormatDisplay   = "02/01/2006"

	// Validation constraints
	MaxReportMonths = 12

	// Transaction types
	TransactionTypeOpeningBalance   = "Saldo Awal"             // openingBalance = sum(unit_link.Ip) * NAV.amount
	TransactionTypeEmployerContrib  = "Iuran Pemberi Kerja"    // Employer amount investment item
	TransactionTypeEmployeeContrib  = "Iuran Peserta"          // Employee amount investment item
	TransactionTypeVoluntaryContrib = "Iuran Sukarela Peserta" // Voluntary amount investment item
	TransactionTypeAdminFee         = "Biaya Admin Iuran"
	TransactionTypeOperationalFee   = "Biaya Pengelolaan Investasi"
)

// TransactionOrder defines the order of transactions in reports
var TransactionOrder = []string{
	TransactionTypeEmployerContrib,
	TransactionTypeEmployeeContrib,
	TransactionTypeVoluntaryContrib,
	TransactionTypeAdminFee,
	TransactionTypeOperationalFee,
}
