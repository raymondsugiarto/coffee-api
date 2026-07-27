package entity

import (
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

// Reference-table values that the upstream system posts
// mutations for. Listed as constants so callers don't drift on
// spelling ("stock_session" vs "stock_session ").
const (
	AccountMutationRefTableStockSession = "stock_session"
	AccountMutationRefTableOrder        = "order"
)

// Reference-module values, grouping upstream sources by domain.
// Reports pivot by module so operators can build a P&L view.
const (
	AccountMutationRefModuleOrder        = "ORDER"
	AccountMutationRefModuleStockSession = "STOCK_SESSION"
)

// AccountMutationDto is the wire shape for one ledger row.
// `amount` is signed: positive = debit to the account, negative =
// credit. The interpretation (asset/expense/etc.) lives at the
// report layer, not here.
type AccountMutationDto struct {
	ID             string  `json:"id"`
	OrganizationID string  `json:"-"`
	AccountID      string  `json:"accountId"     validate:"required"`
	Amount         float64 `json:"amount"        validate:"required"`
	Description    string  `json:"description"`
	RefID          string  `json:"refId"         validate:"required"`
	RefTable       string  `json:"refTable"      validate:"required,oneof=stock_session order"`
	RefModule      string  `json:"refModule"     validate:"required,oneof=ORDER STOCK_SESSION"`
}

func NewAccountMutationDtoFromModel(m *model.AccountMutation) *AccountMutationDto {
	if m == nil {
		return nil
	}
	return &AccountMutationDto{
		ID:             m.ID,
		OrganizationID: m.OrganizationID,
		AccountID:      m.AccountID,
		Amount:         m.Amount,
		Description:    m.Description,
		RefID:          m.RefID,
		RefTable:       m.RefTable,
		RefModule:      m.RefModule,
	}
}

func (d *AccountMutationDto) ToModel() *model.AccountMutation {
	m := &model.AccountMutation{
		OrganizationID: d.OrganizationID,
		AccountID:      d.AccountID,
		Amount:         d.Amount,
		Description:    d.Description,
		RefID:          d.RefID,
		RefTable:       d.RefTable,
		RefModule:      d.RefModule,
	}
	if d.ID != "" {
		m.ID = d.ID
	}
	return m
}

// AccountMutationFindAllRequest powers GET /api/account-mutations.
// Filters supported: by account, by upstream ref, by module. Use
// From/To on `created_at` (YYYY-MM-DD bound) for period reports.
//
// RefID is optional and joins on the exact upstream row when set;
// the listing uses it together with RefTable so a caller can
// pull every mutation that posted against a single stock_session.
type AccountMutationFindAllRequest struct {
	FindAllRequest
	AccountID string
	RefID     string
	RefTable  string
	RefModule string
	// From / To: YYYY-MM-DD strings, inclusive bounds on the
	// stored created_at timestamp. Empty = unbounded.
	From string
	To   string
}

func (r *AccountMutationFindAllRequest) GenerateFilter() {
	if r.AccountID != "" {
		r.Filter = append(r.Filter, pagination.FilterItem{Field: "account_id", Op: "eq", Val: r.AccountID})
	}
	if r.RefID != "" {
		r.Filter = append(r.Filter, pagination.FilterItem{Field: "ref_id", Op: "eq", Val: r.RefID})
	}
	if r.RefTable != "" {
		r.Filter = append(r.Filter, pagination.FilterItem{Field: "ref_table", Op: "eq", Val: r.RefTable})
	}
	if r.RefModule != "" {
		r.Filter = append(r.Filter, pagination.FilterItem{Field: "ref_module", Op: "eq", Val: r.RefModule})
	}
	if r.From != "" {
		r.Filter = append(r.Filter, pagination.FilterItem{Field: "created_at", Op: "gte", Val: r.From})
	}
	if r.To != "" {
		r.Filter = append(r.Filter, pagination.FilterItem{Field: "created_at", Op: "lte", Val: r.To})
	}
}
