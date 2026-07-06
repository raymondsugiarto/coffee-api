package model

import "github.com/raymondsugiarto/coffee-api/pkg/model/concern"

type CashAdjustment struct {
	concern.CommonWithIDs
	SessionID string
	Session   *StockSession
	Type      string // SHORTAGE | OVERAGE
	Amount    float64
	Reason    string
}
