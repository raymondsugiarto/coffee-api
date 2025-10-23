package model

import (
	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
)

type BenefitParticipation struct {
	concern.CommonWithIDs

	OrganizationID string
	Organization   Organization
	CustomerID     string
	Customer       Customer
	ParticipantID  string
	Participant    Participant
	Status         BenefitParticipationStatus

	ExternalDplkName                *string
	ExternalDplkParticipantNumber   *string
	ExternalDplkMonthlyContribution *float64
	HasBpjsPensionProgram           *bool
	InvestmentID                    string

	Details []*BenefitParticipationDetail
}

type BenefitParticipationStatus string

const (
	BenefitParticipationStatusPending  BenefitParticipationStatus = "PENDING"
	BenefitParticipationStatusActive   BenefitParticipationStatus = "ACTIVE"
	BenefitParticipationStatusInactive BenefitParticipationStatus = "INACTIVE"
	BenefitParticipationStatusRejected BenefitParticipationStatus = "REJECTED"
)
