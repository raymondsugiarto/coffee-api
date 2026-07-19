package model

import "github.com/raymondsugiarto/coffee-api/pkg/model/concern"

type SalaryComponent struct {
	concern.CommonWithIDs
	OrganizationID string
	CompanyID      string
	Company        *Company
	ComponentType  string
	MinimumTarget  float64
	Amount         float64
}
