package model

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
)

type Order struct {
	concern.CommonWithIDs
	OrganizationID string
	Organization   *Organization
	CompanyID      string
	Company        *Company
	AdminID        string
	Admin          *Admin
	CustomerID     string
	Code           string
	OrderAt        time.Time
	TotalQty       int
	TotalAmount    float64
	Status         string
	OrderItems     []OrderItem
	OrderPayments  []OrderPayment
}
