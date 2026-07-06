package model

import "github.com/raymondsugiarto/coffee-api/pkg/model/concern"

type PaymentDetail struct {
	concern.CommonWithIDs
	SessionID       string
	Session         *StockSession
	PaymentMethod   string // CASH | QRIS | TRANSFER | OTHER
	Amount          float64
	ReferenceNumber string
	Notes           string
}
