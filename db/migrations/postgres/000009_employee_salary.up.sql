-- ============================================================
-- 000009: employee_salary + employee_salary_component
-- ============================================================
-- employee_salary captures one payroll run for a single employee
-- over a closed date range. The header row stores the totals; the
-- child rows in employee_salary_component store the per-session
-- breakdown that rolled into those totals.
--
-- component_type values mirror the existing salary_component enum
-- plus a new COMMISSION bucket:
--   MEAL_ALLOWANCE
--   ATTENDANCE          (kept for back-compat; always 0 going forward
--                        because bonus-hadir was removed in 000008)
--   COMMISSION          (per-stock-session commission total)
--   BONUS_TARGET
--
-- ref_id / ref_table / ref_source let the system trace each
-- component row back to its origin. Today the only source is
-- SALES (a closed stock_session), but the columns are kept
-- generic so future ref-types (BONUS_MANUAL, ADJUSTMENT, ...) can
-- land without a schema change.

CREATE TABLE IF NOT EXISTS employee_salary (
    id                          varchar(255) PRIMARY KEY,
    organization_id             varchar(255),
    admin_id_employee           varchar(255) NOT NULL,
    start_date                  date         NOT NULL,
    end_date                    date         NOT NULL,
    total_meal_allowance        numeric(20, 4) NOT NULL DEFAULT 0,
    total_attendance_allowance  numeric(20, 4) NOT NULL DEFAULT 0,
    total_commission            numeric(20, 4) NOT NULL DEFAULT 0,
    total_bonus_target          numeric(20, 4) NOT NULL DEFAULT 0,
    total_salary                numeric(20, 4) NOT NULL DEFAULT 0,
    total_cash_receipt          numeric(20, 4) NOT NULL DEFAULT 0,
    remaining_salary            numeric(20, 4) NOT NULL DEFAULT 0,
    created_at                  TIMESTAMP     NOT NULL,
    updated_at                  TIMESTAMP     NULL,
    deleted_at                  TIMESTAMP     NULL
);

CREATE INDEX IF NOT EXISTS idx_employee_salary_admin_date
    ON employee_salary (admin_id_employee, start_date, end_date);

CREATE TABLE IF NOT EXISTS employee_salary_component (
    id                  varchar(255) PRIMARY KEY,
    employee_salary_id  varchar(255) NOT NULL,
    component_type      varchar(40)  NOT NULL,
    amount              numeric(20, 4) NOT NULL DEFAULT 0,
    ref_id              varchar(255),
    ref_table           varchar(60),
    ref_source          varchar(40)  NOT NULL DEFAULT 'SALES',
    created_at          TIMESTAMP     NOT NULL,
    updated_at          TIMESTAMP     NULL,
    deleted_at          TIMESTAMP     NULL
);

CREATE INDEX IF NOT EXISTS idx_employee_salary_component_header
    ON employee_salary_component (employee_salary_id);

CREATE INDEX IF NOT EXISTS idx_employee_salary_component_ref
    ON employee_salary_component (ref_table, ref_id);