package model

import "github.com/raymondsugiarto/coffee-api/pkg/model/concern"

type StockSessionItem struct {
	concern.CommonWithIDs
	SessionID            string
	Session              *StockSession
	ItemID               string
	Item                 *Item
	OutQty               int
	ReturnQty            int
	SoldQty              int
	SellingPriceSnapshot float64
	CostPriceSnapshot    float64
	Subtotal             float64
}
