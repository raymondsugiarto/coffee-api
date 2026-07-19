-- ============================================================
-- 000008: per-session salary breakdown columns on stock_session
-- ============================================================
-- Each row stores the resolved amount the driver is owed for the
-- corresponding salary component:
--   meal_allowance  : flat allowance (typically one row per company
--                     with minimum_target = 0).
--   attendance      : bonus for hitting a minimum attendance target.
--   bonus_target    : bonus for hitting a minimum sales target.
--   total_salary    : sum of the three above. Computed server-side
--                     on close (and recomputed on every write) so
--                     reports don't need to re-derive it.
--
-- Defaults to 0 so existing sessions stay consistent until the
-- service runs its first compute pass.

ALTER TABLE stock_session
    ADD COLUMN IF NOT EXISTS meal_allowance numeric(20, 4) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS attendance     numeric(20, 4) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS bonus_target   numeric(20, 4) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS total_salary    numeric(20, 4) NOT NULL DEFAULT 0;