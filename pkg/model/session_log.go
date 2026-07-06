package model

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
)

// SessionLog is append-only audit log for stock sessions
type SessionLog struct {
	concern.CommonWithIDs
	SessionID string
	Session   *StockSession
	Action    string
	AdminID   string
	Detail    string
	CreatedAt time.Time
}
