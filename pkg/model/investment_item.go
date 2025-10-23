package model

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
)

type InvestmentItem struct {
	concern.CommonWithIDs
	Code                string
	OrganizationID      string
	Organization        *Organization
	InvestmentID        string
	Investment          *Investment
	ParticipantID       string
	Participant         *Participant
	CustomerID          string
	Customer            *Customer
	Type                InvestmentType // DKP or PPIP
	InvestmentType      InvestmentType
	InvestmentProductID string
	InvestmentProduct   *InvestmentProduct
	Percent             float64
	Amount              float64
	FeeAmount           float64
	TotalAmount         float64
	EmployerAmount      float64
	EmployeeAmount      float64
	VoluntaryAmount     float64
	EducationFundAmount float64
	ExpiredAt           time.Time
	InvestmentAt        time.Time
	Status              InvestmentStatus
}
