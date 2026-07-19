-- ============================================================
-- 000009: employee_salary + employee_salary_component (down)
-- ============================================================
DROP INDEX IF EXISTS idx_employee_salary_component_ref;
DROP INDEX IF EXISTS idx_employee_salary_component_header;
DROP INDEX IF EXISTS idx_employee_salary_admin_date;
DROP TABLE IF EXISTS employee_salary_component;
DROP TABLE IF EXISTS employee_salary;