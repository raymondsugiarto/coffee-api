-- ============================================================
-- 000007: commission on stock_session_item + total_commission on
-- stock_session
-- ============================================================
-- Per-row driver commission. commission_snapshot is the per-unit
-- rate captured at write time (sourced from item.commision);
-- commission_total is rate × sold_qty for that row.
-- stock_session.total_commission is the rolled-up sum across all
-- rows. Both default to 0 so existing rows are not affected.

ALTER TABLE stock_session_item
    ADD COLUMN IF NOT EXISTS commission_snapshot numeric(20, 4) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS commission_total    numeric(20, 4) NOT NULL DEFAULT 0;

ALTER TABLE stock_session
    ADD COLUMN IF NOT EXISTS total_commission numeric(20, 4) NOT NULL DEFAULT 0;