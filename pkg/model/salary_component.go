package model

import "github.com/raymondsugiarto/coffee-api/pkg/model/concern"

type ComponentType string

const (
	ComponentTypeMealAllowance ComponentType = "MEAL_ALLOWANCE"
	ComponentTypeCommission    ComponentType = "COMMISSION"
	ComponentTypeBonusTarget   ComponentType = "BONUS_TARGET"
	ComponentTypeAttendance    ComponentType = "ATTENDANCE"
)

type SalaryComponent struct {
	concern.CommonWithIDs
	OrganizationID string
	CompanyID      string
	Company        *Company
	ComponentType  ComponentType
	MinimumTarget  float64
	Amount         float64
}
