package model

import "github.com/raymondsugiarto/coffee-api/pkg/model/concern"

type InvestmentProduct struct {
	concern.CommonWithIDs
	OrganizationID string
	Organization   *Organization
	Code           string
	Name           string
	Description    string
	FixIncome      float64
	StockValue     float64
	MixedValue     float64
	MoneyMarket    float64
	ShariaValue    float64
	FundFactSheet  string
	Riplay         string
	AdminFee       float64
	ManagementFee  float64
	FounderFee     float64
	CommissionFee  float64
	Status         string
	OrderingNumber int
}
