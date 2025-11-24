package model

import "github.com/raymondsugiarto/coffee-api/pkg/model/concern"

type Item struct {
	concern.CommonWithIDs
	OrganizationID string
	Organization   *Organization
	Code           string
	Name           string
	Price          float64
	ItemCompany    []ItemCompany
}
