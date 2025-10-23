package model

import (
	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
)

type CustomerPoint struct {
	concern.CommonWithIDs
	OrganizationID string
	Organization   *Organization
	CustomerID     string
	Customer       *Customer
	Point          float64
	Direction      string
	Description    string
	RefID          string
	RefCode        string
	RefModule      string
}
