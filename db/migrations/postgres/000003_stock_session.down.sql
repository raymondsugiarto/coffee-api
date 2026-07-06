DROP TABLE IF EXISTS session_log;
DROP TABLE IF EXISTS cash_adjustment;
DROP TABLE IF EXISTS payment_detail;
DROP TABLE IF EXISTS stock_session_item;
DROP TABLE IF EXISTS stock_session;
DROP TABLE IF EXISTS item_category;

-- Revert item table changes
ALTER TABLE item
    DROP COLUMN IF EXISTS is_active,
    DROP COLUMN IF EXISTS cost_price,
    DROP COLUMN IF EXISTS category_id,
    DROP COLUMN IF EXISTS sku;
DROP INDEX IF EXISTS idx_item_active;
DROP INDEX IF EXISTS idx_item_category;
DROP INDEX IF EXISTS idx_item_sku;

-- Remove seeded drivers
DELETE FROM admin WHERE id IN ('drv-001-budi', 'drv-002-andi', 'drv-003-siti');
