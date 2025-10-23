package model

import (
	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
)

type Portfolio struct {
	concern.CommonWithIDs
	OrganizationID      string
	Organization        *Organization
	CompanyID           string
	Company             *Company
	CustomerID          string
	Customer            *Customer
	ParticipantID       string
	Participant         *Participant
	InvestmentProductID string
	InvestmentProduct   *InvestmentProduct
	Ip                  float64
}
