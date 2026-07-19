package model

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
)

// EmployeeSalary is the header row for one payroll run for a
// single employee over a closed date range. The components that
// rolled into these totals live in EmployeeSalaryComponent.
type EmployeeSalary struct {
	concern.CommonWithIDs
	OrganizationID     string
	AdminIDEmployee    string
	Admin              *Admin
	StartDate          time.Time
	EndDate            time.Time
	TotalMealAllowance float64
	TotalAttendance    float64
	TotalCommission    float64
	TotalBonusTarget   float64
	TotalSalary        float64
	TotalCashReceipt   float64
	RemainingSalary    float64
}

// EmployeeSalaryComponent is one line on the payroll breakdown.
// component_type mirrors the existing salary_component enum plus
// COMMISSION:
//
//	MEAL_ALLOWANCE | ATTENDANCE | COMMISSION | BONUS_TARGET
//
// ref_id + ref_table + ref_source let the system trace each
// component row back to its origin (today: stock_session; future:
// manual bonus, adjustment, ...).
type EmployeeSalaryComponent struct {
	concern.CommonWithIDs
	EmployeeSalaryID string
	ComponentType    string
	Amount           float64
	RefID            string
	RefTable         string
	RefSource        string // SALES | MANUAL | ADJUSTMENT | ...
}
