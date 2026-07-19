-- Split per-row `sold_qty` into two columns so operators can
-- capture exactly how the driver sold each item, not only the sum.
--
--   cash_sold_qty      — units sold for cash
--   cashless_sold_qty  — units sold via QRIS / cashless
--
-- `sold_qty` stays as a denormalised `cash_sold_qty + cashless_sold_qty`
-- for everything downstream that already reads it (commissions,
-- subtotals, reports). The CHECK constraint is widened so the
-- backend and DB agree on non-negativity.
--
-- Backfill: every existing row is split — all sales are treated as
-- cash sales. That matches the legacy operator UI (single Terjual
-- input) so historical data carries forward without distortion.

ALTER TABLE stock_session_item
    ADD COLUMN IF NOT EXISTS cash_sold_qty INT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS cashless_sold_qty INT NOT NULL DEFAULT 0;

-- Backfill: legacy sold_qty → cash_sold_qty, cashless = 0.
UPDATE stock_session_item
   SET cash_sold_qty = sold_qty
 WHERE cash_sold_qty = 0 AND cashless_sold_qty = 0 AND sold_qty > 0;

-- Soften the non-negative check to include the two new columns.
ALTER TABLE stock_session_item
    DROP CONSTRAINT IF EXISTS chk_ssi_qty_nonneg;
ALTER TABLE stock_session_item
    ADD CONSTRAINT chk_ssi_qty_nonneg
    CHECK (
        out_qty >= 0
        AND return_qty >= 0
        AND sold_qty >= 0
        AND cash_sold_qty >= 0
        AND cashless_sold_qty >= 0
    );
