package entity

import "github.com/raymondsugiarto/coffee-api/pkg/model"

type OrderItemInputDto struct {
	ItemID string
	Item   *ItemDto
	Qty    int
	Price  float64
}

func (i *OrderItemInputDto) ToDto() *OrderItemDto {
	d := &OrderItemDto{
		ItemID:   i.ItemID,
		Item:     i.Item,
		Qty:      i.Qty,
		Price:    i.Price,
		Subtotal: float64(i.Qty) * i.Price,
	}
	return d
}

type OrderItemDto struct {
	ItemID   string
	Item     *ItemDto
	Qty      int
	Price    float64
	Subtotal float64
}

func NewOrderItemDtoFromModel(m *model.OrderItem) *OrderItemDto {
	if m == nil {
		return nil
	}
	return &OrderItemDto{
		ItemID:   m.ItemID,
		Qty:      m.Qty,
		Price:    m.Price,
		Subtotal: m.Subtotal,
	}
}

func (d *OrderItemDto) ToModel() *model.OrderItem {
	m := &model.OrderItem{
		ItemID:   d.ItemID,
		Qty:      d.Qty,
		Price:    d.Price,
		Subtotal: d.Subtotal,
	}
	if d.Item != nil {
		m.Item = d.Item.ToModel()
	}
	return m
}

type OrderItemPerItemCountDto struct {
	ItemName   string  `json:"itemName"`
	TotalQty   int     `json:"totalQty"`
	TotalPrice float64 `json:"totalPrice"`
}
