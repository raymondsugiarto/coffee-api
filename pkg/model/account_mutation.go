package model

import (
	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
)

// AccountMutation is an append-only ledger row recording one
// financial movement against an Account. The pair
// (ref_table, ref_id) links the row back to the upstream entity
// (e.g. "stock_session"/<sessionId>) so reports can reconcile
// without re-deriving amounts.
//
// `ref_module` groups upstream sources (ORDER / STOCK_SESSION
// / ...) so reports can pivot by module. It's intentionally a
// free-form VARCHAR validated at write time by the service layer
// rather than a hard enum because the module list will grow over
// time and we don't want each new module to require a migration.
type AccountMutation struct {
	concern.CommonWithIDs
	OrganizationID string
	AccountID      string
	Amount         float64
	Description    string
	RefID          string
	RefTable       string
	RefModule      string
	// No Account pointer here — ref counted from the FK only.
	// Reports join via the column instead of an association to
	// keep writes fast and predictable.
}
