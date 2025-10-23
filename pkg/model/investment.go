package model

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
)

type Investment struct {
	concern.CommonWithIDs
	OrganizationID     string
	Organization       *Organization
	CompanyID          string
	Company            *Company
	ParticipantID      string
	Participant        *Participant
	CustomerID         string
	Customer           *Customer
	Code               string
	Type               InvestmentType
	Amount             float64
	ExpiredAt          time.Time
	InvestmentAt       time.Time
	Status             InvestmentStatus
	InvestmentPayments []*InvestmentPayment
	InvestmentItems    []*InvestmentItem
	Source             InvestmentSource
}

type InvestmentSource string

const (
	InvestmentSourceRegular              InvestmentSource = "REGULAR"
	InvestmentSourceBenefitParticipation InvestmentSource = "BENEFIT_PARTICIPATION"
)

type InvestmentStatus string

const (
	InvestmentStatusCreated InvestmentStatus = "CREATED"
	InvestmentStatusRequest InvestmentStatus = "REQUEST"
	InvestmentStatusSuccess InvestmentStatus = "SUCCESS"
)

type InvestmentType string

const (
	InvestmentTypeDKP  InvestmentType = "DKP"
	InvestmentTypePPIP InvestmentType = "PPIP"
)
