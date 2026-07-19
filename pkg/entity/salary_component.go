package entity

import (
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

// Salary component types — wire enum. Stored as VARCHAR in the
// `salary_component` table and validated against these constants.
const (
	SalaryComponentTypeMealAllowance = "MEAL_ALLOWANCE"
	SalaryComponentTypeAttendance    = "ATTENDANCE"
	SalaryComponentTypeBonusTarget   = "BONUS_TARGET"
)

type SalaryComponentDto struct {
	ID             string  `json:"id"`
	OrganizationID string  `json:"-"`
	CompanyID      string  `json:"companyId" validate:"required"`
	ComponentType  string  `json:"componentType" validate:"required,oneof=MEAL_ALLOWANCE ATTENDANCE BONUS_TARGET"`
	MinimumTarget  float64 `json:"minimumTarget" validate:"gte=0"`
	Amount         float64 `json:"amount" validate:"gte=0"`
}

func NewSalaryComponentDtoFromModel(m *model.SalaryComponent) *SalaryComponentDto {
	if m == nil {
		return nil
	}
	return &SalaryComponentDto{
		ID:             m.ID,
		OrganizationID: m.OrganizationID,
		CompanyID:      m.CompanyID,
		ComponentType:  m.ComponentType,
		MinimumTarget:  m.MinimumTarget,
		Amount:         m.Amount,
	}
}

func (d *SalaryComponentDto) ToModel() *model.SalaryComponent {
	m := &model.SalaryComponent{
		OrganizationID: d.OrganizationID,
		CompanyID:      d.CompanyID,
		ComponentType:  d.ComponentType,
		MinimumTarget:  d.MinimumTarget,
		Amount:         d.Amount,
	}
	if d.ID != "" {
		m.ID = d.ID
	}
	return m
}

type SalaryComponentFindAllRequest struct {
	FindAllRequest
	CompanyID     string
	ComponentType string
}

func (r *SalaryComponentFindAllRequest) GenerateFilter() {
	if r.CompanyID != "" {
		r.Filter = append(r.Filter, pagination.FilterItem{Field: "company_id", Op: "eq", Val: r.CompanyID})
	}
	if r.ComponentType != "" {
		r.Filter = append(r.Filter, pagination.FilterItem{Field: "component_type", Op: "eq", Val: r.ComponentType})
	}
}
