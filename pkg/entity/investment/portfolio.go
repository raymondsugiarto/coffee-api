package entity

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
)

type Portfolio struct {
	ID                  string                `json:"id"`
	OrganizationID      string                `json:"organizationId,omitempty"`
	CompanyID           string                `json:"companyId,omitempty"`
	Company             *entity.CompanyDto    `json:"company,omitempty"`
	CustomerID          string                `json:"customerId,omitempty"`
	Customer            *entity.CustomerDto   `json:"customer,omitempty"`
	InvestmentProductID string                `json:"investmentProductId,omitempty"`
	InvestmentProduct   *InvestmentProductDto `json:"investmentProduct,omitempty"`
	Amount              float64               `json:"amount,omitempty"`
	NetAssetValueAmount float64               `json:"netAssetValueAmount,omitempty"`
	PortfolioDate       time.Time             `json:"portfolioDate,omitempty"`
}
