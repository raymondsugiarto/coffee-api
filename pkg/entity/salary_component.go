package entity

import (
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

// ComponentType is the entity-side typed enum for the salary
// component_type column. The wire shape carries the raw string,
// and the JSON marshaler renders the underlying value (same as
// model.ComponentType), so the typed field is transparent over
// the API boundary.
//
// Mirrors the entity.UserType convention: a typed string with
// constants re-declared on the entity side. Services switch on
// the entity constants (see entity.StockSessionDto.RecomputeSalary),
// and the model ↔ entity boundary casts via ToModel / FromModel.
type ComponentType string

const (
	ComponentTypeMealAllowance ComponentType = "MEAL_ALLOWANCE"
	ComponentTypeAttendance    ComponentType = "ATTENDANCE"
	ComponentTypeBonusTarget   ComponentType = "BONUS_TARGET"
	// ComponentTypeCommission is intentionally NOT in the DTO
	// validator below: commission is computed per stock-session,
	// not maintained as a master salary_component row. The
	// constant exists so service-layer switches (see
	// entity.StockSessionDto.RecomputeSalary) can match against
	// the typed value without string-cast acrobatics.
	ComponentTypeCommission ComponentType = "COMMISSION"
)

type SalaryComponentDto struct {
	ID             string        `json:"id"`
	OrganizationID string        `json:"-"`
	CompanyID      string        `json:"companyId" validate:"required"`
	ComponentType  ComponentType `json:"componentType" validate:"required,oneof=MEAL_ALLOWANCE ATTENDANCE BONUS_TARGET"`
	MinimumTarget  float64       `json:"minimumTarget" validate:"gte=0"`
	Amount         float64       `json:"amount" validate:"gte=0"`
}

func NewSalaryComponentDtoFromModel(m *model.SalaryComponent) *SalaryComponentDto {
	if m == nil {
		return nil
	}
	return &SalaryComponentDto{
		ID:             m.ID,
		OrganizationID: m.OrganizationID,
		CompanyID:      m.CompanyID,
		ComponentType:  ComponentType(m.ComponentType),
		MinimumTarget:  m.MinimumTarget,
		Amount:         m.Amount,
	}
}

func (d *SalaryComponentDto) ToModel() *model.SalaryComponent {
	m := &model.SalaryComponent{
		OrganizationID: d.OrganizationID,
		CompanyID:      d.CompanyID,
		ComponentType:  model.ComponentType(d.ComponentType),
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
	ComponentType ComponentType
}

func (r *SalaryComponentFindAllRequest) GenerateFilter() {
	if r.CompanyID != "" {
		r.Filter = append(r.Filter, pagination.FilterItem{Field: "company_id", Op: "eq", Val: r.CompanyID})
	}
	if r.ComponentType != "" {
		r.Filter = append(r.Filter, pagination.FilterItem{Field: "component_type", Op: "eq", Val: string(r.ComponentType)})
	}
}
