package entity

import "github.com/raymondsugiarto/coffee-api/pkg/model"

type OrderPaymentDto struct {
	PaymentMethodCode string
}

func NewOrderPaymentDtoFromModel(m *model.OrderPayment) *OrderPaymentDto {
	if m == nil {
		return nil
	}
	return &OrderPaymentDto{
		PaymentMethodCode: m.PaymentMethodCode,
	}
}

func (d *OrderPaymentDto) ToModel() *model.OrderPayment {
	m := &model.OrderPayment{
		PaymentMethodCode: d.PaymentMethodCode,
	}
	return m
}
