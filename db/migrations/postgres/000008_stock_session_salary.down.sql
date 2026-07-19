-- ============================================================
-- 000008: stock_session salary columns (down)
-- ============================================================
ALTER TABLE stock_session
    DROP COLUMN IF EXISTS meal_allowance,
    DROP COLUMN IF EXISTS attendance,
    DROP COLUMN IF EXISTS bonus_target,
    DROP COLUMN IF EXISTS total_salary;