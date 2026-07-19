-- ============================================================
-- 000011: cash_debt column on stock_session
-- ============================================================
-- cash_debt is the operator-entered amount of money the driver
-- owes the company at close (e.g. when an on-the-road expense
-- was paid out of the till, or when the float needs to be
-- topped-up from the previous session's deposit). It is the
-- inverse of cash_receipt and likewise defaults to 0 so every
-- existing row stays consistent until a close writes a value.
--
-- Use it together with cash_receipt against total_cash to
-- reconcile what should be in the till:
--
--   expected_cash = total_cash - cash_debt
--
-- The column is additive so future ops can layer more fields
-- (cash_adjustments, expenses, etc.) without a schema break.

ALTER TABLE stock_session
    ADD COLUMN IF NOT EXISTS cash_debt numeric(20, 4) NOT NULL DEFAULT 0;