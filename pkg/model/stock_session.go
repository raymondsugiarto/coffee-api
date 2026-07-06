package model

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
)

type StockSession struct {
	concern.CommonWithIDs
	OrganizationID string
	EmployeeID     string
	Employee       *Admin
	Date           time.Time
	Status         string // OPEN | CLOSED
	OpenedAt       time.Time
	ClosedAt       *time.Time
	TotalSales     float64
	TotalCash      float64
	TotalQris      float64
	TotalOther     float64
	TotalPayment   float64
	Difference     float64
	TotalItems     int
	Notes          string
	CreatedBy      string
	Items          []StockSessionItem `gorm:"foreignKey:SessionID;constraint:OnDelete:CASCADE"`
	Payments       []PaymentDetail    `gorm:"foreignKey:SessionID;constraint:OnDelete:CASCADE"`
	Adjustments    []CashAdjustment   `gorm:"foreignKey:SessionID;constraint:OnDelete:CASCADE"`
	Logs           []SessionLog       `gorm:"foreignKey:SessionID;constraint:OnDelete:CASCADE"`
}
