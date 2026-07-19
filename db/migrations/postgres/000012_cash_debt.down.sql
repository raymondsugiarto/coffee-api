-- ============================================================
-- 000012: cash_debt table (down)
-- ============================================================
DROP INDEX IF EXISTS idx_cash_debt_org_date;
DROP INDEX IF EXISTS idx_cash_debt_employee_date;
DROP TABLE IF EXISTS cash_debt;