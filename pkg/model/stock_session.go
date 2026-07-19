package model

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
)

type StockSession struct {
	concern.CommonWithIDs
	OrganizationID  string
	EmployeeID      string
	Employee        *Admin
	Date            time.Time
	Status          string // OPEN | CLOSED
	OpenedAt        time.Time
	ClosedAt        *time.Time
	TotalSales      float64
	TotalCash       float64
	TotalQris       float64
	TotalOther      float64
	TotalPayment    float64
	Difference      float64
	TotalItems      int
	TotalCommission float64
	// Salary breakdown resolved from salary_component for the
	// employee's company. Computed on close (and on every write)
	// so reports don't need to re-derive it.
	MealAllowance float64
	Attendance    float64
	BonusTarget   float64
	TotalSalary   float64
	// CashDebt is the operator-entered amount of money the
	// driver owes the company at close (paid-out expenses,
	// top-up float, etc.). It is the inverse of the cash
	// the driver collects — reconciliation rule:
	//
	//   expected_cash = total_cash - cash_debt
	CashDebt    float64
	Notes       string
	CreatedBy   string
	Items       []StockSessionItem `gorm:"foreignKey:SessionID;constraint:OnDelete:CASCADE"`
	Payments    []PaymentDetail    `gorm:"foreignKey:SessionID;constraint:OnDelete:CASCADE"`
	Adjustments []CashAdjustment   `gorm:"foreignKey:SessionID;constraint:OnDelete:CASCADE"`
	Logs        []SessionLog       `gorm:"foreignKey:SessionID;constraint:OnDelete:CASCADE"`
}
