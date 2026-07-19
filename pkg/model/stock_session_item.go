package model

import "github.com/raymondsugiarto/coffee-api/pkg/model/concern"

type StockSessionItem struct {
	concern.CommonWithIDs
	SessionID string
	Session   *StockSession
	ItemID    string
	Item      *Item
	OutQty    int
	ReturnQty int
	// SoldQty is the rolled-up total (cash + cashless). Kept as a
	// denormalised column for downstream reads that already consume
	// it (commission, subtotal, reports). At write time the service
	// always stores it as `CashSoldQty + CashlessSoldQty`.
	SoldQty int
	// CashSoldQty + CashlessSoldQty break the sale into the two
	// sides of the close form. Operators capture them separately
	// during stock-session close so we can audit cash vs QRIS mix
	// per item, not only at the payment-detail layer.
	CashSoldQty          int
	CashlessSoldQty      int
	SellingPriceSnapshot float64
	CostPriceSnapshot    float64
	Subtotal             float64
	// CommissionSnapshot is the driver per-unit commission rate
	// captured at session write time. Recomputed in RecomputeTotals
	// from item.commision × soldQty.
	CommissionSnapshot float64
	// CommissionTotal is the row-level commission amount
	// (commissionSnapshot × soldQty). Summed into
	// stock_session.total_commission.
	CommissionTotal float64
}
