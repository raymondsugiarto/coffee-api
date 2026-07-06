package model

import "github.com/raymondsugiarto/coffee-api/pkg/model/concern"

type Item struct {
	concern.CommonWithIDs
	OrganizationID string
	Organization   *Organization
	CategoryID     string
	Category       *ItemCategory
	Code           string
	SKU            string
	Name           string
	Price          float64
	CostPrice      float64
	IsActive       bool
	ItemCompany    []ItemCompany
}

type ItemCategory struct {
	concern.CommonWithIDs
	OrganizationID string
	Name           string
}
