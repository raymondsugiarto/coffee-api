-- Reverse 000013: collapse cash + cashless into sold_qty, drop
-- the split columns, restore the original CHECK constraint.

-- Soften the new columns' non-negative check (originally these
-- columns don't exist, so we drop the add-on as we go back).
ALTER TABLE stock_session_item
    DROP CONSTRAINT IF EXISTS chk_ssi_qty_nonneg;

-- Roll sold_qty up from the two split columns before dropping them
-- so we don't lose the historical data.
UPDATE stock_session_item
   SET sold_qty = cash_sold_qty + cashless_sold_qty;

ALTER TABLE stock_session_item
    DROP COLUMN IF EXISTS cash_sold_qty,
    DROP COLUMN IF EXISTS cashless_sold_qty;

-- Restore original constraint (taken from migration 000003).
ALTER TABLE stock_session_item
    ADD CONSTRAINT chk_ssi_qty_nonneg
    CHECK (out_qty >= 0 AND return_qty >= 0 AND sold_qty >= 0);
