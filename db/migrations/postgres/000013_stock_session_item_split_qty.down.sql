
ALTER TABLE stock_session_item
    DROP COLUMN IF EXISTS cash_sold_qty,
    DROP COLUMN IF EXISTS cashless_sold_qty;
