package model

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
)

// CashDebt is a driver-issued cash advance / petty-cash entry.
// One row per advance. Settlement happens out-of-band when the
// driver hands back receipts or the company deducts from the
// daily till — so this ledger stays purely additive.
//
// We deliberately keep `AdminIDEmployee` as a plain string and
// do NOT define a `*Admin` relation. GORM requires a `foreignKey`
// tag whenever a relation field is present, and there is no FK
// at the DB level (the employee column is just an admin-id
// reference for filtering). The entity layer resolves the
// employee separately when needed.
type CashDebt struct {
	concern.CommonWithIDs
	OrganizationID  string
	AdminIDEmployee string
	Date            time.Time
	Amount          float64
	PaymentMethod   string // CASH | CASHLESS
	Notes           string
}
