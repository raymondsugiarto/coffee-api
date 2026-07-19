-- ============================================================
-- 000006: salary_component master (down)
-- ============================================================
DROP INDEX IF EXISTS idx_salary_component_org;
DROP INDEX IF EXISTS idx_salary_component_type;
DROP INDEX IF EXISTS idx_salary_component_company;
DROP TABLE IF EXISTS salary_component;