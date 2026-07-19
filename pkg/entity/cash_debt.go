package entity

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

// CashDebt payment methods — wire enum.
const (
	CashDebtPaymentMethodCash     = "CASH"
	CashDebtPaymentMethodCashless = "CASHLESS"
)

type CashDebtDto struct {
	ID              string  `json:"id"`
	OrganizationID  string  `json:"-"`
	AdminIDEmployee string  `json:"adminIdEmployee" validate:"required"`
	Date            string  `json:"date" validate:"required,len=10"` // YYYY-MM-DD
	Amount          float64 `json:"amount" validate:"gte=0"`
	PaymentMethod   string  `json:"paymentMethod" validate:"required,oneof=CASH CASHLESS"`
	Notes           string  `json:"notes"`
}

func NewCashDebtDtoFromModel(m *model.CashDebt) *CashDebtDto {
	if m == nil {
		return nil
	}
	return &CashDebtDto{
		ID:              m.ID,
		OrganizationID:  m.OrganizationID,
		AdminIDEmployee: m.AdminIDEmployee,
		Date:            m.Date.Format("2006-01-02"),
		Amount:          m.Amount,
		PaymentMethod:   m.PaymentMethod,
		Notes:           m.Notes,
	}
}

func (d *CashDebtDto) ToModel() *model.CashDebt {
	parsedDate, _ := time.Parse("2006-01-02", d.Date)
	m := &model.CashDebt{
		OrganizationID:  d.OrganizationID,
		AdminIDEmployee: d.AdminIDEmployee,
		Date:            parsedDate,
		Amount:          d.Amount,
		PaymentMethod:   d.PaymentMethod,
		Notes:           d.Notes,
	}
	if d.ID != "" {
		m.ID = d.ID
	}
	return m
}

// CashDebtFindAllRequest powers GET /api/cash-debts.
type CashDebtFindAllRequest struct {
	FindAllRequest
	AdminIDEmployee string
	From            string // YYYY-MM-DD
	To              string // YYYY-MM-DD
	PaymentMethod   string
}

func (r *CashDebtFindAllRequest) GenerateFilter() {
	if r.AdminIDEmployee != "" {
		r.Filter = append(r.Filter, pagination.FilterItem{Field: "admin_id_employee", Op: "eq", Val: r.AdminIDEmployee})
	}
	if r.PaymentMethod != "" {
		r.Filter = append(r.Filter, pagination.FilterItem{Field: "payment_method", Op: "eq", Val: r.PaymentMethod})
	}
	if r.From != "" {
		r.Filter = append(r.Filter, pagination.FilterItem{Field: "date", Op: "gte", Val: r.From})
	}
	if r.To != "" {
		r.Filter = append(r.Filter, pagination.FilterItem{Field: "date", Op: "lte", Val: r.To})
	}
}
