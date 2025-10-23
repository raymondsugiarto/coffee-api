package model

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
)

type UnitLink struct {
	concern.CommonWithIDs
	OrganizationID      string
	Organization        Organization
	ParticipantID       string
	Participant         *Participant
	CustomerID          string
	Customer            *Customer
	InvestmentProductID string
	InvestmentProduct   *InvestmentProduct
	Type                InvestmentType
	TransactionDate     time.Time
	TotalAmount         float64
	Nab                 float64
	Ip                  float64
	CreatedAt           time.Time
}
