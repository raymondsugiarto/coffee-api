package model

import "github.com/raymondsugiarto/coffee-api/pkg/model/concern"

type OrderItem struct {
	concern.CommonWithIDs
	OrderID  string
	Order    *Order
	ItemID   string
	Item     *Item
	Qty      int
	Price    float64
	Subtotal float64
}
