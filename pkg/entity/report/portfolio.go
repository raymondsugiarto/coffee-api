package entity

import (
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
)

type PortfolioFindAllRequest struct {
	entity.FindAllRequest
	CustomerID    string `json:"customerId"`
	ParticipantID string `json:"participantId"`
}

type PortfolioReportDto struct {
	ID                  string  `json:"id"` // Generated as participant_id-investment_product_id
	ParticipantID       string  `json:"participantId"`
	CustomerID          string  `json:"customerId"`
	InvestmentProductID string  `json:"investmentProductId"`
	Ip                  float64 `json:"ip"` // Total units (SUM of all transactions)
	LatestNav           float64 `json:"latestNav"`
	TotalBalance        float64 `json:"totalBalance"`

	// Data from JOIN for nested objects
	CustomerFirstName     string `json:"-" gorm:"column:customer_first_name"`
	ParticipantCode       string `json:"-" gorm:"column:participant_code"`
	InvestmentProductCode string `json:"-" gorm:"column:investment_product_code"`
	InvestmentProductName string `json:"-" gorm:"column:investment_product_name"`

	// Nested objects for frontend compatibility
	Customer          *CustomerBasicDto          `json:"customer"`
	Participant       *ParticipantBasicDto       `json:"participant"`
	InvestmentProduct *InvestmentProductBasicDto `json:"investmentProduct"`
}

type CustomerBasicDto struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
}

type ParticipantBasicDto struct {
	ID   string `json:"id"`
	Code string `json:"code"`
}

type InvestmentProductBasicDto struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// PopulateNestedObjects fills the nested objects from the JOIN query results
func (p *PortfolioReportDto) PopulateNestedObjects() {
	p.Customer = &CustomerBasicDto{
		ID:        p.CustomerID,
		FirstName: p.CustomerFirstName,
	}

	p.Participant = &ParticipantBasicDto{
		ID:   p.ParticipantID,
		Code: p.ParticipantCode,
	}

	p.InvestmentProduct = &InvestmentProductBasicDto{
		ID:   p.InvestmentProductID,
		Code: p.InvestmentProductCode,
		Name: p.InvestmentProductName,
	}
}
