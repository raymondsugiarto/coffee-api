package entity

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/model"
)

type OrderInputDto struct {
	AdminID       string
	Admin         *AdminDto
	OrderAt       time.Time
	TotalQty      int
	TotalAmount   float64
	PaymentMethod string
	OrderItems    []OrderItemInputDto
}

func (i *OrderInputDto) ToDto() *OrderDto {
	d := &OrderDto{
		AdminID: i.AdminID,
		Admin:   i.Admin,
		OrderAt: i.OrderAt,
	}
	totalQty := 0
	totalAmount := 0.0
	for _, oi := range i.OrderItems {
		d.OrderItems = append(d.OrderItems, *oi.ToDto())
		totalAmount += oi.Price * float64(oi.Qty)
		totalQty = totalQty + oi.Qty
	}
	d.TotalQty = totalQty
	d.TotalAmount = totalAmount
	d.Status = "SUCCESS"
	d.OrderPayments = []OrderPaymentDto{
		{
			PaymentMethodCode: i.PaymentMethod,
		},
	}
	return d
}

type OrderDto struct {
	ID             string            `json:"id"`
	OrganizationID string            `json:"-"`
	Organization   *OrganizationDto  `json:"-"`
	CompanyID      string            `json:"-"`
	Company        *CompanyDto       `json:"-"`
	UserID         string            `json:"-"`
	AdminID        string            `json:"-"`
	Admin          *AdminDto         `json:"-"`
	CustomerID     string            `json:"-"`
	Code           string            `json:"code"`
	OrderAt        time.Time         `json:"orderAt"`
	TotalQty       int               `json:"totalQty"`
	TotalAmount    float64           `json:"totalAmount"`
	Status         string            `json:"status"`
	OrderItems     []OrderItemDto    `json:"orderItems"`
	OrderPayments  []OrderPaymentDto `json:"orderPayments"`
}

func NewOrderDtoFromModel(m *model.Order) *OrderDto {
	if m == nil {
		return nil
	}
	return &OrderDto{
		ID:             m.ID,
		OrganizationID: m.OrganizationID,
		CompanyID:      m.CompanyID,
		AdminID:        m.AdminID,
		CustomerID:     m.CustomerID,
		Code:           m.Code,
		OrderAt:        m.OrderAt,
		TotalQty:       m.TotalQty,
		TotalAmount:    m.TotalAmount,
		Status:         m.Status,
	}
}

func (d *OrderDto) ToModel() *model.Order {
	m := &model.Order{
		OrganizationID: d.OrganizationID,
		CompanyID:      d.CompanyID,
		AdminID:        d.AdminID,
		CustomerID:     d.CustomerID,
		Code:           d.Code,
		OrderAt:        d.OrderAt,
		TotalQty:       d.TotalQty,
		TotalAmount:    d.TotalAmount,
		Status:         d.Status,
	}
	for _, oi := range d.OrderItems {
		m.OrderItems = append(m.OrderItems, *oi.ToModel())
	}
	for _, op := range d.OrderPayments {
		m.OrderPayments = append(m.OrderPayments, *op.ToModel())
	}
	if d.ID != "" {
		m.ID = d.ID
	}
	return m
}

type OrderFindAllRequest struct {
	FindAllRequest
	UserID         string
	AdminID        string
	CompanyID      string
	MyEmployeeItem bool
}

func (r *OrderFindAllRequest) GenerateFilter() {
}
