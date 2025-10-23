package model

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
)

type TransactionFee struct {
	concern.CommonWithIDs
	OrganizationID      string
	Organization        Organization
	CompanyID           string
	Company             *Company
	ParticipantID       string
	Participant         *Participant
	InvestmentProductID string
	InvestmentProduct   *InvestmentProduct
	Type                InvestmentType
	TransactionDate     time.Time
	OperationFee        float64
	Nav                 float64
	Ip                  float64
	PortfolioAmount     float64
}
