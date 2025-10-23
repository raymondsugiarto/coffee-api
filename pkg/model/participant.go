package model

import (
	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
)

type ParticipantStatus string

const (
	ParticipantStatusInactive ParticipantStatus = "INACTIVE"
	ParticipantStatusActive   ParticipantStatus = "ACTIVE"
)

type Participant struct {
	concern.CommonWithIDs
	OrganizationID string
	Organization   *Organization
	CustomerID     string
	Customer       *Customer
	Code           string
	Type           InvestmentType
	Status         ParticipantStatus
}
