-- ============================================================
-- 000005: commision on item
-- ============================================================
-- Commision is the per-unit commission paid to the driver on this
-- product. Stored as numeric so we keep full precision; defaults
-- to 0 so existing rows are not affected.

ALTER TABLE item
    ADD COLUMN IF NOT EXISTS commision numeric(20, 4) NOT NULL DEFAULT 0;