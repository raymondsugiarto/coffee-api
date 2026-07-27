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

// OpenStockSessionItemInputDto is the wire shape for a single item
// when opening a stock session. OutQty is the morning count the
// driver physically takes out (>= 1). ReturnQty is optional (the
// morning flow typically sends 0).
type OpenStockSessionItemInputDto struct {
	ItemID    string   `json:"itemId" validate:"required"`
	Item      *ItemDto `json:"item,omitempty"`
	OutQty    int      `json:"outQty" validate:"required,gte=1"`
	ReturnQty int      `json:"returnQty" validate:"gte=0"`
}

// CloseStockSessionItemInputDto is the wire shape for a single item
// when closing a stock session. The operator breaks the sales into
// two buckets at close time, since the operator's mental model is
// "Terjual (Cash)" + "Terjual (QRIS)":
//
//   - CashSoldQty:      units paid in cash
//   - CashlessSoldQty:  units paid via QRIS / cashless
//   - ReturnQty:        units that came back unsold
//
// The morning OutQty is *not* sent because it is already persisted
// on the session row from open time and does not change at close
// time. The service layer derives the rolled `SoldQty` as
// `CashSoldQty + CashlessSoldQty`, and the persisted OutQty is
// `CashSoldQty + CashlessSoldQty + ReturnQty` revalidated against
// the original morning count.
type CloseStockSessionItemInputDto struct {
	ItemID          string `json:"itemId"                       validate:"required"`
	CashSoldQty     int    `json:"cashSoldQty"                 validate:"gte=0"`
	CashlessSoldQty int    `json:"cashlessSoldQty"             validate:"gte=0"`
	ReturnQty       int    `json:"returnQty"                   validate:"gte=0"`
}

// SoldQty is a derived helper exposed for readers that prefer the
// rolled view (commission, subtotal, frontend display).
func (i *CloseStockSessionItemInputDto) SoldQty() int {
	return i.CashSoldQty + i.CashlessSoldQty
}

// StockSessionItemInputDto is the *internal* shape used inside the
// service layer after the wire payload has been normalised. Both
// OpenStockSessionItemInputDto and CloseStockSessionItemInputDto map
// into this shape via ToStockSessionItemInputDto(). Keep this struct
// out of API request types.
type StockSessionItemInputDto struct {
	ItemID    string   `json:"-"`
	Item      *ItemDto `json:"-"`
	OutQty    int      `json:"-"`
	ReturnQty int      `json:"-"`
	// SoldQty is the rolled value the rest of the system reads.
	// It is always `CashSoldQty + CashlessSoldQty`. Persisted as
	// denormalised sold_qty + the two split columns.
	SoldQty         int `json:"-"`
	CashSoldQty     int `json:"-"`
	CashlessSoldQty int `json:"-"`
}

// ToStockSessionItemInputDto normalises an open wire item into the
// internal input DTO. OutQty is the morning count; soldQty and the
// two split columns default to 0 because the morning flow has no
// notion of cash vs QRIS yet (recording happens at close).
func (i *OpenStockSessionItemInputDto) ToStockSessionItemInputDto() *StockSessionItemInputDto {
	sold := i.OutQty - i.ReturnQty
	if sold < 0 {
		sold = 0
	}
	return &StockSessionItemInputDto{
		ItemID:    i.ItemID,
		Item:      i.Item,
		OutQty:    i.OutQty,
		ReturnQty: i.ReturnQty,
		SoldQty:   sold,
	}
}

// ToStockSessionItemInputDto normalises a close wire item into the
// internal input DTO. SoldQty is reconstructed as cash + cashless;
// the service layer revalidates that against the morning count
// already persisted on the session row.
func (i *CloseStockSessionItemInputDto) ToStockSessionItemInputDto() *StockSessionItemInputDto {
	cash := i.CashSoldQty
	if cash < 0 {
		cash = 0
	}
	cashless := i.CashlessSoldQty
	if cashless < 0 {
		cashless = 0
	}
	sold := cash + cashless
	outQty := sold + i.ReturnQty
	return &StockSessionItemInputDto{
		ItemID:          i.ItemID,
		Item:            nil, // close path: item is never nested on the wire
		OutQty:          outQty,
		ReturnQty:       i.ReturnQty,
		SoldQty:         sold,
		CashSoldQty:     cash,
		CashlessSoldQty: cashless,
	}
}

// ToDto maps the internal input DTO into the persisted shape.
func (i *StockSessionItemInputDto) ToDto() *StockSessionItemDto {
	sold := i.SoldQty
	if sold < 0 {
		sold = 0
	}
	cash := i.CashSoldQty
	if cash < 0 {
		cash = 0
	}
	cashless := i.CashlessSoldQty
	if cashless < 0 {
		cashless = 0
	}
	// If the row didn't carry the split columns (legacy / open path
	// replayed after close), fall back to assigning everything to
	// cash so the persisted columns stay in sync with rolled
	// SoldQty.
	if cash == 0 && cashless == 0 && sold > 0 {
		cash = sold
	}
	outQty := i.OutQty
	if outQty == 0 {
		outQty = sold + i.ReturnQty
	}
	subtotal := float64(sold) * i.SellingPriceSnapshot()
	return &StockSessionItemDto{
		ItemID:               i.ItemID,
		Item:                 i.Item,
		OutQty:               outQty,
		ReturnQty:            i.ReturnQty,
		SoldQty:              sold,
		CashSoldQty:          cash,
		CashlessSoldQty:      cashless,
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
	ID        string   `json:"id"`
	SessionID string   `json:"-"`
	ItemID    string   `json:"itemId"`
	Item      *ItemDto `json:"item,omitempty"`
	OutQty    int      `json:"outQty"`
	ReturnQty int      `json:"returnQty"`
	// SoldQty is the rolled-up total (cash + cashless). Kept for
	// downstream consumers (commission calc, subtotal rollup,
	// reports) that already read this column.
	SoldQty int `json:"soldQty"`
	// CashSoldQty + CashlessSoldQty split the same sale into the
	// two halves at item granularity. They mirror what the close
	// form captures (Terjual (Cash) + Terjual (QRIS)). Always in
	// sync with SoldQty via RecomputeTotals / ToDto.
	CashSoldQty          int     `json:"cashSoldQty"`
	CashlessSoldQty      int     `json:"cashlessSoldQty"`
	SellingPriceSnapshot float64 `json:"sellingPriceSnapshot"`
	CostPriceSnapshot    float64 `json:"costPriceSnapshot"`
	Subtotal             float64 `json:"subtotal"`
	// CommissionSnapshot is the per-unit commission rate at the
	// time of the write; CommissionTotal is rate × soldQty.
	CommissionSnapshot float64 `json:"commissionSnapshot"`
	CommissionTotal    float64 `json:"commissionTotal"`
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
		CashSoldQty:          m.CashSoldQty,
		CashlessSoldQty:      m.CashlessSoldQty,
		SellingPriceSnapshot: m.SellingPriceSnapshot,
		CostPriceSnapshot:    m.CostPriceSnapshot,
		Subtotal:             m.Subtotal,
		CommissionSnapshot:   m.CommissionSnapshot,
		CommissionTotal:      m.CommissionTotal,
	}
	if m.Item != nil {
		d.Item = NewItemDtoFromModel(m.Item)
	}
	return d
}

// ToModel maps a DTO row to its persistence model.
//
// NOTE: this method intentionally does NOT carry over `d.Item` to the
// model. Stock-session open/update only stores a reference (item_id)
// to an existing item; the nested Item DTO is treated as informational
// at most. Mapping it through to the model would cause GORM's
// association auto-save to INSERT into the item table on every
// session write, which is not what we want.
func (d *StockSessionItemDto) ToModel() *model.StockSessionItem {
	m := &model.StockSessionItem{
		SessionID:            d.SessionID,
		ItemID:               d.ItemID,
		OutQty:               d.OutQty,
		ReturnQty:            d.ReturnQty,
		SoldQty:              d.SoldQty,
		CashSoldQty:          d.CashSoldQty,
		CashlessSoldQty:      d.CashlessSoldQty,
		SellingPriceSnapshot: d.SellingPriceSnapshot,
		CostPriceSnapshot:    d.CostPriceSnapshot,
		Subtotal:             d.Subtotal,
		CommissionSnapshot:   d.CommissionSnapshot,
		CommissionTotal:      d.CommissionTotal,
	}
	if d.ID != "" {
		m.ID = d.ID
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
	EmployeeID string                         `json:"employeeId" validate:"required"`
	Date       string                         `json:"date" validate:"required"` // YYYY-MM-DD
	Notes      string                         `json:"notes"`
	Items      []OpenStockSessionItemInputDto `json:"items" validate:"required,min=1,dive"`
}

type CloseStockSessionInputDto struct {
	Items       []CloseStockSessionItemInputDto `json:"items" validate:"required,min=1,dive"`
	Payments    []PaymentDetailInputDto         `json:"payments" validate:"required,min=1,dive"`
	Adjustments []CashAdjustmentInputDto        `json:"adjustments" validate:"omitempty,dive"`
	Notes       string                          `json:"notes"`
}

type StockSessionDto struct {
	ID                  string     `json:"id"`
	OrganizationID      string     `json:"-"`
	EmployeeID          string     `json:"employeeId"`
	Employee            *AdminDto  `json:"employee,omitempty"`
	Date                string     `json:"date"` // YYYY-MM-DD
	Status              string     `json:"status"`
	OpenedAt            time.Time  `json:"openedAt"`
	ClosedAt            *time.Time `json:"closedAt,omitempty"`
	TotalSales          float64    `json:"totalSales"`
	TotalCash           float64    `json:"totalCash"`
	TotalQris           float64    `json:"totalQris"`
	TotalOther          float64    `json:"totalOther"`
	TotalPayment        float64    `json:"totalPayment"`
	Difference          float64    `json:"difference"`
	TotalItems          int        `json:"totalItems"`
	TotalCommission     float64    `json:"totalCommission"`
	MinTargetCommission float64    `json:"minTargetCommission"` // derived from salary_component, used to determine if the driver met the minimum commission target for bonus
	// Salary breakdown resolved from salary_component for the
	// driver's company. Computed server-side on every write.
	MealAllowance float64 `json:"mealAllowance"`
	Attendance    float64 `json:"attendance"`
	BonusTarget   float64 `json:"bonusTarget"`
	TotalSalary   float64 `json:"totalSalary"`
	// CashDebt is the operator-entered amount the driver owes
	// the company at close.
	CashDebt    float64               `json:"cashDebt"`
	Notes       string                `json:"notes"`
	CreatedBy   string                `json:"createdBy"`
	Items       []StockSessionItemDto `json:"items"`
	Payments    []PaymentDetailDto    `json:"payments,omitempty"`
	Adjustments []CashAdjustmentDto   `json:"adjustments,omitempty"`
}

func NewStockSessionDtoFromModel(m *model.StockSession) *StockSessionDto {
	if m == nil {
		return nil
	}
	d := &StockSessionDto{
		ID:              m.ID,
		OrganizationID:  m.OrganizationID,
		EmployeeID:      m.EmployeeID,
		Date:            m.Date.Format("2006-01-02"),
		Status:          m.Status,
		OpenedAt:        m.OpenedAt,
		ClosedAt:        m.ClosedAt,
		TotalSales:      m.TotalSales,
		TotalCash:       m.TotalCash,
		TotalQris:       m.TotalQris,
		TotalOther:      m.TotalOther,
		TotalPayment:    m.TotalPayment,
		Difference:      m.Difference,
		TotalItems:      m.TotalItems,
		TotalCommission: m.TotalCommission,
		MealAllowance:   m.MealAllowance,
		Attendance:      m.Attendance,
		BonusTarget:     m.BonusTarget,
		TotalSalary:     m.TotalSalary,
		CashDebt:        m.CashDebt,
		Notes:           m.Notes,
		CreatedBy:       m.CreatedBy,
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
		OrganizationID:      d.OrganizationID,
		EmployeeID:          d.EmployeeID,
		Status:              d.Status,
		OpenedAt:            d.OpenedAt,
		ClosedAt:            d.ClosedAt,
		TotalSales:          d.TotalSales,
		TotalCash:           d.TotalCash,
		TotalQris:           d.TotalQris,
		TotalOther:          d.TotalOther,
		TotalPayment:        d.TotalPayment,
		Difference:          d.Difference,
		TotalItems:          d.TotalItems,
		TotalCommission:     d.TotalCommission,
		MinTargetCommission: d.MinTargetCommission,
		MealAllowance:       d.MealAllowance,
		Attendance:          d.Attendance,
		BonusTarget:         d.BonusTarget,
		TotalSalary:         d.TotalSalary,
		CashDebt:            d.CashDebt,
		Notes:               d.Notes,
		CreatedBy:           d.CreatedBy,
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
//
// Per-row `soldQty` resolution:
//
//   - Open-time rows and items that were already in the morning
//     session at close time carry `OutQty > 0` (the morning count).
//     We derive `sold = outQty − returnQty` and write it back. The
//     split columns default to `0 + 0`, since the morning flow
//     has no notion of cash vs QRIS — the close form re-records
//     them later.
//   - Rows that did NOT exist in the morning session (variants
//     added on the close form, typically child variants sharing
//     their parent's pool) carry `OutQty == 0`. We can't compute
//     `sold` from that, so we trust the admin's `soldQty`
//     (which equals `cashSoldQty + cashlessSoldQty`) exactly as
//     they typed it.
//
// Subtotal is always `soldQty × price`, regardless of how `soldQty`
// was resolved.
//
// Per-row commission:
//   - CommissionSnapshot = item.commision (captured from the
//     master item table; preloaded on read or hydrated from the
//     item table on open/update/close).
//   - CommissionTotal = commissionSnapshot × soldQty.
//
// `d.TotalCommission` is the sum of every row's CommissionTotal.
func (d *StockSessionDto) RecomputeTotals() {
	var totalSales float64
	var totalCommission float64
	var totalItems int
	for i, it := range d.Items {
		var sold int
		// New (close-only) item: the morning count doesn't
		// apply. Trust the admin's split. Keep them as-typed
		// and derive soldQty from the sum.
		cash := it.CashSoldQty
		if cash < 0 {
			cash = 0
		}
		cashless := it.CashlessSoldQty
		if cashless < 0 {
			cashless = 0
		}
		sold = cash + cashless
		if sold < 0 {
			sold = 0
		}
		d.Items[i].SoldQty = sold
		d.Items[i].CashSoldQty = cash
		d.Items[i].CashlessSoldQty = cashless
		subtotal := float64(sold) * it.SellingPriceSnapshot
		d.Items[i].Subtotal = subtotal
		totalSales += subtotal
		totalItems += sold

		// Pull the commission rate from the master item, falling
		// back to whatever snapshot the row already carries so
		// historical sessions never see their commission reset to
		// 0 if the master item is later edited.
		rate := d.Items[i].CommissionSnapshot
		if d.Items[i].Item != nil {
			rate = d.Items[i].Item.Commision
			d.Items[i].CommissionSnapshot = rate
		}
		commission := float64(sold) * rate
		d.Items[i].CommissionTotal = commission
		totalCommission += commission
	}
	d.TotalSales = totalSales
	d.TotalItems = totalItems
	d.TotalCommission = totalCommission

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
	d.TotalCash = cash - d.CashDebt
	if d.TotalCash < 0 {
		d.TotalCash = 0
	}
	d.TotalQris = qris
	d.TotalOther = other
	// TotalPayment mirrors the same deduction so the reported total
	// is always the net of what the driver actually brings back.
	d.TotalPayment = totalPayment - d.CashDebt
	if d.TotalPayment < 0 {
		d.TotalPayment = 0
	}

	// Difference = Net Payment - Total Sales.
	// Net payment already has CashDebt deducted so Difference
	// reflects what should physically be in the till vs what was
	// expected from sales.
	d.Difference = d.TotalPayment - totalSales
}

// RecomputeSalary resolves the per-session salary breakdown
// (meal_allowance, attendance, bonus_target) from a list of salary
// components that apply to the driver's company.
//
// Resolution rules:
//
//   - meal_allowance : sum of every MEAL_ALLOWANCE row in the
//     company (typically just one row, with minimum_target = 0).
//   - attendance     : pick the ATTENDANCE row with the highest
//     minimum_target that the driver still clears. `totalQty`
//     here is typically the session's TotalItems — the user
//     can refine later to use a separate days-worked metric.
//   - bonus_target   : same shape as attendance, applied to
//     BONUS_TARGET rows.
//
// `totalSalary` is the sum of the three. Components that don't
// match any row resolve to 0 — that's the safe default for new
// companies that haven't set up their salary bands yet.
func (d *StockSessionDto) RecomputeSalary(components []SalaryComponentDto, totalQty int) {
	// Attendance bonus calculation has been removed by product
	// decision. The `Attendance` column is preserved on the model
	// and DTO for historical/reporting parity, but is always 0
	// on every write.
	var meal, bonus, minTargetCommission float64
	var attendance = 0.0
	var bestBonusTarget = -1
	for _, c := range components {
		switch c.ComponentType {
		case ComponentTypeCommission:
			minTargetCommission = c.MinimumTarget * c.Amount
		case ComponentTypeMealAllowance:
			meal += c.Amount
		case ComponentTypeAttendance:
			// Intentionally ignored: bonus-hadir removed.
			// Field kept at 0 for column/UI stability.
			if c.MinimumTarget <= float64(totalQty) {
				attendance = c.Amount
			}
		case ComponentTypeBonusTarget:
			if c.MinimumTarget <= float64(totalQty) &&
				c.MinimumTarget > float64(bestBonusTarget) {
				bestBonusTarget = int(c.MinimumTarget)
				bonus = c.Amount
			}
		}
	}
	d.MealAllowance = meal
	d.Attendance = attendance
	d.BonusTarget = bonus
	d.TotalSalary = (d.TotalCommission - minTargetCommission) + meal + bonus + attendance
	d.MinTargetCommission = minTargetCommission
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
	Date               string                 `json:"date"`
	Sessions           int                    `json:"sessions"`
	TotalSales         float64                `json:"totalSales"`
	TotalCash          float64                `json:"totalCash"`
	TotalQris          float64                `json:"totalQris"`
	TotalOther         float64                `json:"totalOther"`
	TotalPayment       float64                `json:"totalPayment"`
	TotalDiff          float64                `json:"totalDifference"`
	TotalCommission    float64                `json:"totalCommission"`
	TotalMealAllowance float64                `json:"totalMealAllowance"`
	TotalBonusTarget   float64                `json:"totalBonusTarget"`
	TotalSalary        float64                `json:"totalSalary"`
	ByEmployee         []EmployeeReportRowDto `json:"byEmployee"`
}

type MonthlyReportDto struct {
	Year               int                    `json:"year"`
	Month              int                    `json:"month"`
	Sessions           int                    `json:"sessions"`
	TotalSales         float64                `json:"totalSales"`
	TotalCash          float64                `json:"totalCash"`
	TotalQris          float64                `json:"totalQris"`
	TotalDiff          float64                `json:"totalDifference"`
	TotalCommission    float64                `json:"totalCommission"`
	TotalMealAllowance float64                `json:"totalMealAllowance"`
	TotalBonusTarget   float64                `json:"totalBonusTarget"`
	TotalSalary        float64                `json:"totalSalary"`
	Daily              []DailyReportDto       `json:"daily"`
	ByEmployee         []EmployeeReportRowDto `json:"byEmployee"`
}

type EmployeeReportRowDto struct {
	EmployeeID    string  `json:"employeeId"`
	EmployeeName  string  `json:"employeeName"`
	Sessions      int     `json:"sessions"`
	TotalSales    float64 `json:"totalSales"`
	TotalCash     float64 `json:"totalCash"`
	TotalQris     float64 `json:"totalQris"`
	Difference    float64 `json:"difference"`
	Commission    float64 `json:"commission"`
	MealAllowance float64 `json:"mealAllowance"`
	BonusTarget   float64 `json:"bonusTarget"`
	TotalSalary   float64 `json:"totalSalary"`
}

type TopProductRowDto struct {
	ProductID   string  `json:"productId"`
	ProductName string  `json:"productName"`
	SKU         string  `json:"sku"`
	TotalQty    int     `json:"totalQty"`
	TotalSales  float64 `json:"totalSales"`
}

type EmployeePerformanceRowDto struct {
	EmployeeID    string  `json:"employeeId"`
	EmployeeName  string  `json:"employeeName"`
	Sessions      int     `json:"sessions"`
	TotalItems    int     `json:"totalItems"`
	TotalSales    float64 `json:"totalSales"`
	TotalCash     float64 `json:"totalCash"`
	TotalQris     float64 `json:"totalQris"`
	TotalDiff     float64 `json:"totalDifference"`
	Commission    float64 `json:"commission"`
	MealAllowance float64 `json:"mealAllowance"`
	BonusTarget   float64 `json:"bonusTarget"`
	TotalSalary   float64 `json:"totalSalary"`
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
