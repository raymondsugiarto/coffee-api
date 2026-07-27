package model

import (
	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
)

// Account is a chart-of-accounts row. Operators create these as
// named categories ("Kas", "Piutang", "Penjualan", "Beban Gaji")
// and then post account_mutation rows against them from upstream
// flows (order, stock_session, payroll, etc.).
//
// `code` is the short identifier printed on reports and used as
// the seedable key per organization. The (organization_id, code)
// pair is unique so different orgs can use the same codes.
type Account struct {
	concern.CommonWithIDs
	OrganizationID string
	Name           string
	Code           string
}
