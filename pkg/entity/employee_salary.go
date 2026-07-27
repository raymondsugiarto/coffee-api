package entity

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/model"
)

// Component types stored on employee_salary_component.component_type.
// Mirrors salary_component's enum plus COMMISSION.
const (
	EmployeeSalaryComponentTypeMealAllowance = "MEAL_ALLOWANCE"
	EmployeeSalaryComponentTypeAttendance    = "ATTENDANCE"
	EmployeeSalaryComponentTypeCommission    = "COMMISSION"
	EmployeeSalaryComponentTypeBonusTarget   = "BONUS_TARGET"
)

// Ref-source values stored on employee_salary_component.ref_source.
const (
	EmployeeSalaryRefSourceSales = "SALES"
)

// ===== EmployeeSalaryComponentDto =====

type EmployeeSalaryComponentDto struct {
	ID               string  `json:"id"`
	EmployeeSalaryID string  `json:"employeeSalaryId"`
	ComponentType    string  `json:"componentType"`
	Amount           float64 `json:"amount"`
	RefID            string  `json:"refId,omitempty"`
	RefTable         string  `json:"refTable,omitempty"`
	RefSource        string  `json:"refSource"`
}

func NewEmployeeSalaryComponentDtoFromModel(m *model.EmployeeSalaryComponent) *EmployeeSalaryComponentDto {
	if m == nil {
		return nil
	}
	return &EmployeeSalaryComponentDto{
		ID:               m.ID,
		EmployeeSalaryID: m.EmployeeSalaryID,
		ComponentType:    m.ComponentType,
		Amount:           m.Amount,
		RefID:            m.RefID,
		RefTable:         m.RefTable,
		RefSource:        m.RefSource,
	}
}

func (d *EmployeeSalaryComponentDto) ToModel() *model.EmployeeSalaryComponent {
	return &model.EmployeeSalaryComponent{
		EmployeeSalaryID: d.EmployeeSalaryID,
		ComponentType:    d.ComponentType,
		Amount:           d.Amount,
		RefID:            d.RefID,
		RefTable:         d.RefTable,
		RefSource:        d.RefSource,
	}
}

// ===== EmployeeSalaryDto =====

type EmployeeSalaryDto struct {
	ID                 string                       `json:"id"`
	OrganizationID     string                       `json:"-"`
	AdminIDEmployee    string                       `json:"adminIdEmployee"`
	StartDate          string                       `json:"startDate"` // YYYY-MM-DD
	EndDate            string                       `json:"endDate"`
	TotalMealAllowance float64                      `json:"totalMealAllowance"`
	TotalAttendance    float64                      `json:"totalAttendanceAllowance"`
	TotalCommission    float64                      `json:"totalCommission"`
	TotalBonusTarget   float64                      `json:"totalBonusTarget"`
	TotalSalary        float64                      `json:"totalSalary"`
	TotalCashReceipt   float64                      `json:"totalCashReceipt"`
	RemainingSalary    float64                      `json:"remainingSalary"`
	Components         []EmployeeSalaryComponentDto `json:"components,omitempty"`
}

func NewEmployeeSalaryDtoFromModel(m *model.EmployeeSalary) *EmployeeSalaryDto {
	if m == nil {
		return nil
	}
	return &EmployeeSalaryDto{
		ID:                 m.ID,
		OrganizationID:     m.OrganizationID,
		AdminIDEmployee:    m.AdminIDEmployee,
		StartDate:          m.StartDate.Format("2006-01-02"),
		EndDate:            m.EndDate.Format("2006-01-02"),
		TotalMealAllowance: m.TotalMealAllowance,
		TotalAttendance:    m.TotalAttendance,
		TotalCommission:    m.TotalCommission,
		TotalBonusTarget:   m.TotalBonusTarget,
		TotalSalary:        m.TotalSalary,
		TotalCashReceipt:   m.TotalCashReceipt,
		RemainingSalary:    m.RemainingSalary,
	}
}

func (d *EmployeeSalaryDto) ToModel() *model.EmployeeSalary {
	startDate, _ := time.Parse("2006-01-02", d.StartDate)
	endDate, _ := time.Parse("2006-01-02", d.EndDate)
	m := &model.EmployeeSalary{
		OrganizationID:     d.OrganizationID,
		AdminIDEmployee:    d.AdminIDEmployee,
		StartDate:          startDate,
		EndDate:            endDate,
		TotalMealAllowance: d.TotalMealAllowance,
		TotalAttendance:    d.TotalAttendance,
		TotalCommission:    d.TotalCommission,
		TotalBonusTarget:   d.TotalBonusTarget,
		TotalSalary:        d.TotalSalary,
		TotalCashReceipt:   d.TotalCashReceipt,
		RemainingSalary:    d.RemainingSalary,
	}
	if d.ID != "" {
		m.ID = d.ID
	}
	return m
}

// ===== SimulateRequest =====

// SimulatePayrollRequest is the wire shape for POST /api/payroll/simulate.
// The backend reads closed stock_sessions for this employee over the
// given date range, walks each session's resolved commission + salary
// breakdown (already stored on the row from close time), and produces
// a SimulatePayrollResultDto the frontend can show before saving.
type SimulatePayrollRequest struct {
	AdminIDEmployee string `json:"adminIdEmployee" validate:"required"`
	StartDate       string `json:"startDate"        validate:"required,len=10"`
	EndDate         string `json:"endDate"          validate:"required,len=10"`
}

// ===== SimulatePayrollResultDto =====

// SimulatePayrollSessionCashDebtDto is one cash_debt entry tied to a
// session's date. A session can have multiple advances, so the
// session rows carries a slice rather than a single number.
type SimulatePayrollSessionCashDebtDto struct {
	ID            string  `json:"id"`
	Date          string  `json:"date"`
	Amount        float64 `json:"amount"`
	PaymentMethod string  `json:"paymentMethod"`
	Notes         string  `json:"notes"`
}

// SimulatePayrollSessionDto is one row in the simulation breakdown.
// It mirrors the per-session totals already persisted on stock_session
// so the operator can verify what will land in employee_salary.
//
// CashDebts holds every cash_debt row whose date matches this
// session's date — the simulator surfaces them so the operator
// can see exactly which advances are being netted from the
// remaining salary, instead of seeing only the rolled total.
type SimulatePayrollSessionDto struct {
	SessionID           string                              `json:"sessionId"`
	Date                string                              `json:"date"`
	Status              string                              `json:"status"`
	TotalSales          float64                             `json:"totalSales"`
	Attendance          float64                             `json:"attendance"`
	MinTargetCommission float64                             `json:"minTargetCommission"`
	Commission          float64                             `json:"commission"`
	MealAllowance       float64                             `json:"mealAllowance"`
	BonusTarget         float64                             `json:"bonusTarget"`
	TotalSalary         float64                             `json:"totalSalary"`
	CashDebts           []SimulatePayrollSessionCashDebtDto `json:"cashDebts"`
}

// SimulatePayrollResultDto is the full preview the frontend renders
// after a SIMULATE call. It carries both the per-session evidence
// and the rolled-up totals the Save action will persist.
//
// TotalCashDebt is the sum of every cash_debt row in the date
// range for this employee. RemainingSalary already nets it out:
//
//	remaining = totalSalary - totalCashReceipt - totalCashDebt
type SimulatePayrollResultDto struct {
	AdminIDEmployee    string                      `json:"adminIdEmployee"`
	StartDate          string                      `json:"startDate"`
	EndDate            string                      `json:"endDate"`
	Sessions           []SimulatePayrollSessionDto `json:"sessions"`
	TotalMealAllowance float64                     `json:"totalMealAllowance"`
	TotalAttendance    float64                     `json:"totalAttendance"`
	TotalCommission    float64                     `json:"totalCommission"`
	TotalBonusTarget   float64                     `json:"totalBonusTarget"`
	TotalSalary        float64                     `json:"totalSalary"`
	TotalCashReceipt   float64                     `json:"totalCashReceipt"`
	TotalCashDebt      float64                     `json:"totalCashDebt"`
	RemainingSalary    float64                     `json:"remainingSalary"`
	SessionCount       int                         `json:"sessionCount"`
}

// EmployeeSalaryFindAllRequest powers GET /api/payroll.
type EmployeeSalaryFindAllRequest struct {
	FindAllRequest
	AdminIDEmployee string
}

// SavePayrollRequest is the wire shape for POST /api/payroll/save.
// The frontend sends the operator-approved totals + per-component
// breakdown (already grouped by component type + ref_id); the
// service flattens this into employee_salary + employee_salary_component.
type SavePayrollRequest struct {
	AdminIDEmployee    string                       `json:"adminIdEmployee" validate:"required"`
	StartDate          string                       `json:"startDate"        validate:"required,len=10"`
	EndDate            string                       `json:"endDate"          validate:"required,len=10"`
	TotalMealAllowance float64                      `json:"totalMealAllowance"`
	TotalAttendance    float64                      `json:"totalAttendanceAllowance"`
	TotalCommission    float64                      `json:"totalCommission"`
	TotalBonusTarget   float64                      `json:"totalBonusTarget"`
	TotalSalary        float64                      `json:"totalSalary"`
	TotalCashReceipt   float64                      `json:"totalCashReceipt"`
	Components         []EmployeeSalaryComponentDto `json:"components"`
}
