-- ============================================================
-- 000010: salary_component + admin_company seed fix (down)
-- ============================================================
-- This down is a no-op: the seed corrections in 000010.up
-- patched existing rows back to canonical IDs. Reverting
-- those patches would risk re-introducing the orphaned
-- 'sekian-company-id' literal that 000010.up fixed.
SELECT 1;