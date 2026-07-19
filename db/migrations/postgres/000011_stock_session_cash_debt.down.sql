-- ============================================================
-- 000011: cash_debt column on stock_session (down)
-- ============================================================
ALTER TABLE stock_session
    DROP COLUMN IF EXISTS cash_debt;