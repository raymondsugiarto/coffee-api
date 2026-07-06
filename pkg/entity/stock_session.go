package entity

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

const (
	StockSessionStatusOpen   = "OPEN"
	StockSessionStatusClosed = "CLOSED"

	PaymentMethodCash     = "CASH"
	PaymentMethodQris     = "QRIS"
	PaymentMethodTransfer = "TRANSFER"
	PaymentMethodOther    = "OTHER"

	CashAdjustmentShortage = "SHORTAGE"
	CashAdjustmentOverage  = "OVERAGE"

	SessionActionOpen   = "OPEN"
	SessionActionUpdate = "UPDATE"
	SessionActionClose  = "CLOSE"
)

// ===== StockSessionItem =====

type StockSessionItemInputDto struct {
	ItemID    string   `json:"itemId" validate:"required"`
	Item      *ItemDto `json:"item,omitempty"`
	OutQty    int      `json:"outQty" validate:"required,gte=1"`
	ReturnQty int      `json:"returnQty" validate:"gte=0"`
}

func (i *StockSessionItemInputDto) ToDto() *StockSessionItemDto {
	sold := i.OutQty - i.ReturnQty
	if sold < 0 {
		sold = 0
	}
	subtotal := float64(sold) * i.SellingPriceSnapshot()
	return &StockSessionItemDto{
		ItemID:               i.ItemID,
		Item:                 i.Item,
		OutQty:               i.OutQty,
		ReturnQty:            i.ReturnQty,
		SoldQty:              sold,
		SellingPriceSnapshot: i.SellingPriceSnapshot(),
		CostPriceSnapshot:    i.CostPriceSnapshot(),
		Subtotal:             subtotal,
	}
}

// SellingPriceSnapshot falls back to item's selling price when not given.
func (i *StockSessionItemInputDto) SellingPriceSnapshot() float64 {
	if i.Item != nil {
		return i.Item.Price
	}
	return 0
}

func (i *StockSessionItemInputDto) CostPriceSnapshot() float64 {
	if i.Item != nil {
		return i.Item.CostPrice
	}
	return 0
}

type StockSessionItemDto struct {
	ID                   string   `json:"id"`
	SessionID            string   `json:"-"`
	ItemID               string   `json:"itemId"`
	Item                 *ItemDto `json:"item,omitempty"`
	OutQty               int      `json:"outQty"`
	ReturnQty            int      `json:"returnQty"`
	SoldQty              int      `json:"soldQty"`
	SellingPriceSnapshot float64  `json:"sellingPriceSnapshot"`
	CostPriceSnapshot    float64  `json:"costPriceSnapshot"`
	Subtotal             float64  `json:"subtotal"`
}

func NewStockSessionItemDtoFromModel(m *model.StockSessionItem) *StockSessionItemDto {
	if m == nil {
		return nil
	}
	d := &StockSessionItemDto{
		ID:                   m.ID,
		SessionID:            m.SessionID,
		ItemID:               m.ItemID,
		OutQty:               m.OutQty,
		ReturnQty:            m.ReturnQty,
		SoldQty:              m.SoldQty,
		SellingPriceSnapshot: m.SellingPriceSnapshot,
		CostPriceSnapshot:    m.CostPriceSnapshot,
		Subtotal:             m.Subtotal,
	}
	if m.Item != nil {
		d.Item = NewItemDtoFromModel(m.Item)
	}
	return d
}

func (d *StockSessionItemDto) ToModel() *model.StockSessionItem {
	m := &model.StockSessionItem{
		SessionID:            d.SessionID,
		ItemID:               d.ItemID,
		OutQty:               d.OutQty,
		ReturnQty:            d.ReturnQty,
		SoldQty:              d.SoldQty,
		SellingPriceSnapshot: d.SellingPriceSnapshot,
		CostPriceSnapshot:    d.CostPriceSnapshot,
		Subtotal:             d.Subtotal,
	}
	if d.ID != "" {
		m.ID = d.ID
	}
	if d.Item != nil {
		m.Item = d.Item.ToModel()
	}
	return m
}

// ===== PaymentDetail =====

type PaymentDetailInputDto struct {
	PaymentMethod   string  `json:"paymentMethod" validate:"required,oneof=CASH QRIS TRANSFER OTHER"`
	Amount          float64 `json:"amount" validate:"gte=0"`
	ReferenceNumber string  `json:"referenceNumber"`
	Notes           string  `json:"notes"`
}

func (i *PaymentDetailInputDto) ToDto() *PaymentDetailDto {
	return &PaymentDetailDto{
		PaymentMethod:   i.PaymentMethod,
		Amount:          i.Amount,
		ReferenceNumber: i.ReferenceNumber,
		Notes:           i.Notes,
	}
}

type PaymentDetailDto struct {
	ID              string  `json:"id"`
	SessionID       string  `json:"-"`
	PaymentMethod   string  `json:"paymentMethod"`
	Amount          float64 `json:"amount"`
	ReferenceNumber string  `json:"referenceNumber"`
	Notes           string  `json:"notes"`
}

func NewPaymentDetailDtoFromModel(m *model.PaymentDetail) *PaymentDetailDto {
	if m == nil {
		return nil
	}
	return &PaymentDetailDto{
		ID:              m.ID,
		SessionID:       m.SessionID,
		PaymentMethod:   m.PaymentMethod,
		Amount:          m.Amount,
		ReferenceNumber: m.ReferenceNumber,
		Notes:           m.Notes,
	}
}

func (d *PaymentDetailDto) ToModel() *model.PaymentDetail {
	m := &model.PaymentDetail{
		SessionID:       d.SessionID,
		PaymentMethod:   d.PaymentMethod,
		Amount:          d.Amount,
		ReferenceNumber: d.ReferenceNumber,
		Notes:           d.Notes,
	}
	if d.ID != "" {
		m.ID = d.ID
	}
	return m
}

// ===== CashAdjustment =====

type CashAdjustmentInputDto struct {
	Type   string  `json:"type" validate:"required,oneof=SHORTAGE OVERAGE"`
	Amount float64 `json:"amount" validate:"gt=0"`
	Reason string  `json:"reason"`
}

func (i *CashAdjustmentInputDto) ToDto() *CashAdjustmentDto {
	return &CashAdjustmentDto{
		Type:   i.Type,
		Amount: i.Amount,
		Reason: i.Reason,
	}
}

type CashAdjustmentDto struct {
	ID        string  `json:"id"`
	SessionID string  `json:"-"`
	Type      string  `json:"type"`
	Amount    float64 `json:"amount"`
	Reason    string  `json:"reason"`
}

func NewCashAdjustmentDtoFromModel(m *model.CashAdjustment) *CashAdjustmentDto {
	if m == nil {
		return nil
	}
	return &CashAdjustmentDto{
		ID:        m.ID,
		SessionID: m.SessionID,
		Type:      m.Type,
		Amount:    m.Amount,
		Reason:    m.Reason,
	}
}

func (d *CashAdjustmentDto) ToModel() *model.CashAdjustment {
	m := &model.CashAdjustment{
		SessionID: d.SessionID,
		Type:      d.Type,
		Amount:    d.Amount,
		Reason:    d.Reason,
	}
	if d.ID != "" {
		m.ID = d.ID
	}
	return m
}

// ===== StockSession =====

type OpenStockSessionInputDto struct {
	EmployeeID string                     `json:"employeeId" validate:"required"`
	Date       string                     `json:"date" validate:"required"` // YYYY-MM-DD
	Notes      string                     `json:"notes"`
	Items      []StockSessionItemInputDto `json:"items" validate:"required,min=1,dive"`
}

type CloseStockSessionInputDto struct {
	Items       []StockSessionItemInputDto `json:"items" validate:"required,min=1,dive"`
	Payments    []PaymentDetailInputDto    `json:"payments" validate:"required,min=1,dive"`
	Adjustments []CashAdjustmentInputDto   `json:"adjustments" validate:"omitempty,dive"`
	Notes       string                     `json:"notes"`
}

type StockSessionDto struct {
	ID             string                `json:"id"`
	OrganizationID string                `json:"-"`
	EmployeeID     string                `json:"employeeId"`
	Employee       *AdminDto             `json:"employee,omitempty"`
	Date           string                `json:"date"` // YYYY-MM-DD
	Status         string                `json:"status"`
	OpenedAt       time.Time             `json:"openedAt"`
	ClosedAt       *time.Time            `json:"closedAt,omitempty"`
	TotalSales     float64               `json:"totalSales"`
	TotalCash      float64               `json:"totalCash"`
	TotalQris      float64               `json:"totalQris"`
	TotalOther     float64               `json:"totalOther"`
	TotalPayment   float64               `json:"totalPayment"`
	Difference     float64               `json:"difference"`
	TotalItems     int                   `json:"totalItems"`
	Notes          string                `json:"notes"`
	CreatedBy      string                `json:"createdBy"`
	Items          []StockSessionItemDto `json:"items"`
	Payments       []PaymentDetailDto    `json:"payments,omitempty"`
	Adjustments    []CashAdjustmentDto   `json:"adjustments,omitempty"`
}

func NewStockSessionDtoFromModel(m *model.StockSession) *StockSessionDto {
	if m == nil {
		return nil
	}
	d := &StockSessionDto{
		ID:             m.ID,
		OrganizationID: m.OrganizationID,
		EmployeeID:     m.EmployeeID,
		Date:           m.Date.Format("2006-01-02"),
		Status:         m.Status,
		OpenedAt:       m.OpenedAt,
		ClosedAt:       m.ClosedAt,
		TotalSales:     m.TotalSales,
		TotalCash:      m.TotalCash,
		TotalQris:      m.TotalQris,
		TotalOther:     m.TotalOther,
		TotalPayment:   m.TotalPayment,
		Difference:     m.Difference,
		TotalItems:     m.TotalItems,
		Notes:          m.Notes,
		CreatedBy:      m.CreatedBy,
	}
	if m.Employee != nil {
		d.Employee = (&AdminDto{}).FromModel(m.Employee)
	}
	for _, it := range m.Items {
		d.Items = append(d.Items, *NewStockSessionItemDtoFromModel(&it))
	}
	for _, p := range m.Payments {
		d.Payments = append(d.Payments, *NewPaymentDetailDtoFromModel(&p))
	}
	for _, a := range m.Adjustments {
		d.Adjustments = append(d.Adjustments, *NewCashAdjustmentDtoFromModel(&a))
	}
	return d
}

func (d *StockSessionDto) ToModel() *model.StockSession {
	m := &model.StockSession{
		OrganizationID: d.OrganizationID,
		EmployeeID:     d.EmployeeID,
		Status:         d.Status,
		OpenedAt:       d.OpenedAt,
		ClosedAt:       d.ClosedAt,
		TotalSales:     d.TotalSales,
		TotalCash:      d.TotalCash,
		TotalQris:      d.TotalQris,
		TotalOther:     d.TotalOther,
		TotalPayment:   d.TotalPayment,
		Difference:     d.Difference,
		TotalItems:     d.TotalItems,
		Notes:          d.Notes,
		CreatedBy:      d.CreatedBy,
	}
	if d.ID != "" {
		m.ID = d.ID
	}
	if d.Date != "" {
		if t, err := time.Parse("2006-01-02", d.Date); err == nil {
			m.Date = t
		}
	}
	for _, it := range d.Items {
		m.Items = append(m.Items, *it.ToModel())
	}
	for _, p := range d.Payments {
		m.Payments = append(m.Payments, *p.ToModel())
	}
	for _, a := range d.Adjustments {
		m.Adjustments = append(m.Adjustments, *a.ToModel())
	}
	return m
}

// RecomputeTotals recalculates TotalSales, payments, and difference.
// Call before Save/Close.
func (d *StockSessionDto) RecomputeTotals() {
	var totalSales float64
	var totalItems int
	for i, it := range d.Items {
		sold := it.OutQty - it.ReturnQty
		if sold < 0 {
			sold = 0
		}
		subtotal := float64(sold) * it.SellingPriceSnapshot
		d.Items[i].SoldQty = sold
		d.Items[i].Subtotal = subtotal
		totalSales += subtotal
		totalItems += sold
	}
	d.TotalSales = totalSales
	d.TotalItems = totalItems

	var cash, qris, other, totalPayment float64
	for _, p := range d.Payments {
		totalPayment += p.Amount
		switch p.PaymentMethod {
		case PaymentMethodCash:
			cash += p.Amount
		case PaymentMethodQris:
			qris += p.Amount
		case PaymentMethodTransfer, PaymentMethodOther:
			other += p.Amount
		}
	}
	d.TotalCash = cash
	d.TotalQris = qris
	d.TotalOther = other
	d.TotalPayment = totalPayment

	// Difference = Total Payment - Total Sales
	// Adjustments are informational and don't affect difference.
	d.Difference = totalPayment - totalSales
}

type StockSessionFindAllRequest struct {
	FindAllRequest
	EmployeeID string
	Date       string // YYYY-MM-DD
	Status     string
	From       string
	To         string
}

func (r *StockSessionFindAllRequest) GenerateFilter() {
	if r.EmployeeID != "" {
		r.Filter = append(r.Filter, pagination.FilterItem{Field: "employee_id", Op: "eq", Val: r.EmployeeID})
	}
	if r.Status != "" {
		r.Filter = append(r.Filter, pagination.FilterItem{Field: "status", Op: "eq", Val: r.Status})
	}
	if r.Date != "" {
		r.Filter = append(r.Filter, pagination.FilterItem{Field: "date", Op: "eq", Val: r.Date})
	}
	if r.From != "" {
		r.Filter = append(r.Filter, pagination.FilterItem{Field: "date", Op: "gte", Val: r.From})
	}
	if r.To != "" {
		r.Filter = append(r.Filter, pagination.FilterItem{Field: "date", Op: "lte", Val: r.To})
	}
}

// ===== SessionLog =====

type SessionLogDto struct {
	ID        string    `json:"id"`
	SessionID string    `json:"sessionId"`
	Action    string    `json:"action"`
	AdminID   string    `json:"adminId"`
	Detail    string    `json:"detail"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewSessionLogDtoFromModel(m *model.SessionLog) *SessionLogDto {
	if m == nil {
		return nil
	}
	return &SessionLogDto{
		ID:        m.ID,
		SessionID: m.SessionID,
		Action:    m.Action,
		AdminID:   m.AdminID,
		Detail:    m.Detail,
		CreatedAt: m.CreatedAt,
	}
}

// ===== Reports =====

type DailyReportDto struct {
	Date         string                 `json:"date"`
	Sessions     int                    `json:"sessions"`
	TotalSales   float64                `json:"totalSales"`
	TotalCash    float64                `json:"totalCash"`
	TotalQris    float64                `json:"totalQris"`
	TotalOther   float64                `json:"totalOther"`
	TotalPayment float64                `json:"totalPayment"`
	TotalDiff    float64                `json:"totalDifference"`
	ByEmployee   []EmployeeReportRowDto `json:"byEmployee"`
}

type MonthlyReportDto struct {
	Year       int                    `json:"year"`
	Month      int                    `json:"month"`
	Sessions   int                    `json:"sessions"`
	TotalSales float64                `json:"totalSales"`
	TotalCash  float64                `json:"totalCash"`
	TotalQris  float64                `json:"totalQris"`
	TotalDiff  float64                `json:"totalDifference"`
	Daily      []DailyReportDto       `json:"daily"`
	ByEmployee []EmployeeReportRowDto `json:"byEmployee"`
}

type EmployeeReportRowDto struct {
	EmployeeID   string  `json:"employeeId"`
	EmployeeName string  `json:"employeeName"`
	Sessions     int     `json:"sessions"`
	TotalSales   float64 `json:"totalSales"`
	TotalCash    float64 `json:"totalCash"`
	TotalQris    float64 `json:"totalQris"`
	Difference   float64 `json:"difference"`
}

type TopProductRowDto struct {
	ProductID   string  `json:"productId"`
	ProductName string  `json:"productName"`
	SKU         string  `json:"sku"`
	TotalQty    int     `json:"totalQty"`
	TotalSales  float64 `json:"totalSales"`
}

type EmployeePerformanceRowDto struct {
	EmployeeID   string  `json:"employeeId"`
	EmployeeName string  `json:"employeeName"`
	Sessions     int     `json:"sessions"`
	TotalItems   int     `json:"totalItems"`
	TotalSales   float64 `json:"totalSales"`
	TotalCash    float64 `json:"totalCash"`
	TotalQris    float64 `json:"totalQris"`
	TotalDiff    float64 `json:"totalDifference"`
}

type DashboardSummaryDto struct {
	TodaySales        float64 `json:"todaySales"`
	TodayCash         float64 `json:"todayCash"`
	TodayQris         float64 `json:"todayQris"`
	TodayTransactions int     `json:"todayTransactions"`
	OpenSessions      int     `json:"openSessions"`
	ClosedSessions    int     `json:"closedSessions"`
	TotalSessions     int     `json:"totalSessions"`
}
