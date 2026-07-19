-- ============================================================
-- 000006: salary_component master
-- ============================================================
-- Driver salary components: meal allowance, attendance bonus,
-- bonus target tier, etc. Per company; each row is keyed on
-- (company_id, component_type, minimum_target). Component types
-- are stored as plain VARCHAR; the Go side enforces the enum
-- (MEAL_ALLOWANCE, ATTENDANCE, BONUS_TARGET).

CREATE TABLE IF NOT EXISTS salary_component (
    id varchar(255) PRIMARY KEY,
    organization_id varchar(255) NULL,
    company_id varchar(255) NOT NULL,
    component_type varchar(64) NOT NULL,
    minimum_target numeric(20, 4) NOT NULL DEFAULT 0,
    amount numeric(20, 4) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX IF NOT EXISTS idx_salary_component_company ON salary_component(company_id);
CREATE INDEX IF NOT EXISTS idx_salary_component_type ON salary_component(component_type);
CREATE INDEX IF NOT EXISTS idx_salary_component_org ON salary_component(organization_id);

-- ============================================================
-- Seed data — two demo companies ("SEKIAN" and "Mawaru")
-- ============================================================
INSERT INTO salary_component (id, organization_id, company_id, component_type, minimum_target, amount, created_at)
VALUES
    -- SEKIAN — meal allowance
    ('sc-sekian-meal',           '16f31f95-b356-4e96-b0df-c7f5052beb95', 'sekian-company-id', 'MEAL_ALLOWANCE', 0,  15000, NOW()),
    -- SEKIAN — attendance bonus (minimum days worked)
    ('sc-sekian-att-27',         '16f31f95-b356-4e96-b0df-c7f5052beb95', 'sekian-company-id', 'ATTENDANCE',    27, 35000, NOW()),
    ('sc-sekian-att-20',         '16f31f95-b356-4e96-b0df-c7f5052beb95', 'sekian-company-id', 'ATTENDANCE',    20, 35000, NOW()),
    -- SEKIAN — bonus target tier (per sales-target hit)
    ('sc-sekian-bonus-30',       '16f31f95-b356-4e96-b0df-c7f5052beb95', 'sekian-company-id', 'BONUS_TARGET',  30, 10000, NOW()),
    ('sc-sekian-bonus-35',       '16f31f95-b356-4e96-b0df-c7f5052beb95', 'sekian-company-id', 'BONUS_TARGET',  35, 15000, NOW()),
    ('sc-sekian-bonus-40',       '16f31f95-b356-4e96-b0df-c7f5052beb95', 'sekian-company-id', 'BONUS_TARGET',  40, 20000, NOW()),

    -- Mawaru — meal allowance
    ('sc-mawaru-meal',           '16f31f95-b356-4e96-b0df-c7f5052beb95', 'mawaru-company-id', 'MEAL_ALLOWANCE', 0,  15000, NOW()),
    -- Mawaru — attendance bonus
    ('sc-mawaru-att-20',         '16f31f95-b356-4e96-b0df-c7f5052beb95', 'mawaru-company-id', 'ATTENDANCE',    20, 35000, NOW()),
    -- Mawaru — bonus target tier
    ('sc-mawaru-bonus-25',       '16f31f95-b356-4e96-b0df-c7f5052beb95', 'mawaru-company-id', 'BONUS_TARGET',  25, 10000, NOW()),
    ('sc-mawaru-bonus-30',       '16f31f95-b356-4e96-b0df-c7f5052beb95', 'mawaru-company-id', 'BONUS_TARGET',  30, 15000, NOW()),
    ('sc-mawaru-bonus-35',       '16f31f95-b356-4e96-b0df-c7f5052beb95', 'mawaru-company-id', 'BONUS_TARGET',  35, 20000, NOW())
ON CONFLICT (id) DO UPDATE
SET component_type  = EXCLUDED.component_type,
    minimum_target  = EXCLUDED.minimum_target,
    amount          = EXCLUDED.amount,
    company_id      = EXCLUDED.company_id,
    organization_id = EXCLUDED.organization_id;