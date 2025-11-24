package model

import "github.com/raymondsugiarto/coffee-api/pkg/model/concern"

type ItemCompany struct {
	concern.CommonWithIDs
	OrganizationID string
	Organization   *Organization
	CompanyID      string
	Company        *Company
	ItemID         string
	Item           *Item
}
