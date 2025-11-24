package model

import "github.com/raymondsugiarto/coffee-api/pkg/model/concern"

type OrderPayment struct {
	concern.CommonWithIDs
	OrderID           string
	Order             *Order
	PaymentMethodCode string
}
