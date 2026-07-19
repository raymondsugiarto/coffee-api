-- ============================================================
-- 000007: stock_session commission columns (down)
-- ============================================================
ALTER TABLE stock_session_item
    DROP COLUMN IF EXISTS commission_snapshot,
    DROP COLUMN IF EXISTS commission_total;

ALTER TABLE stock_session
    DROP COLUMN IF EXISTS total_commission;