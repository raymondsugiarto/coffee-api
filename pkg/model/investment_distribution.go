package model

import (
	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
)

type InvestmentDistribution struct {
	concern.CommonWithIDs
	OrganizationID        string
	Organization          *Organization
	Type                  string
	CompanyID             string
	Company               *Company
	ParticipantID         string
	Participant           *Participant
	CustomerID            string
	Customer              *Customer
	InvestmentProductID   string
	InvestmentProduct     *InvestmentProduct
	Percent               float64
	BaseContribution      float64
	VoluntaryContribution float64
}
