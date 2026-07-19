-- ============================================================
-- 000010: fix salary_component + admin_company seed data
-- ============================================================
-- 000006_salary_component seeded salary_component rows with
-- `company_id` literals ('sekian-company-id', 'mawaru-company-id')
-- that don't exist in the `company` table. The real company
-- IDs (from db/migrations/seed.sql) are UUIDs:
--   Sekian = a93150c9-eb99-4c62-8fc1-c414c8a0f78d
--   Mawaru = 3b80cdb0-5a52-42e5-9a9b-38f7c3e9164f
--
-- This migration also binds the seeded test drivers
-- ('drv-001-budi', 'drv-002-andi', 'drv-003-siti' from
-- 000003_stock_session) to Sekian via admin_company, so the
-- payroll simulate endpoint can resolve their company and
-- load the matching salary_component rows.

UPDATE salary_component
SET company_id = CASE company_id
    WHEN 'sekian-company-id' THEN 'a93150c9-eb99-4c62-8fc1-c414c8a0f78d'
    WHEN 'mawaru-company-id' THEN '3b80cdb0-5a52-42e5-9a9b-38f7c3e9164f'
    ELSE company_id
END
WHERE company_id IN ('sekian-company-id', 'mawaru-company-id');

INSERT INTO admin_company (id, organization_id, company_id, admin_id, created_at, updated_at)
VALUES
    ('ac-drv-001-budi',  '16f31f95-b356-4e96-b0df-c7f5052beb95', 'a93150c9-eb99-4c62-8fc1-c414c8a0f78d', 'drv-001-budi',  NOW(), NOW()),
    ('ac-drv-002-andi',  '16f31f95-b356-4e96-b0df-c7f5052beb95', 'a93150c9-eb99-4c62-8fc1-c414c8a0f78d', 'drv-002-andi',  NOW(), NOW()),
    ('ac-drv-003-siti',  '16f31f95-b356-4e96-b0df-c7f5052beb95', 'a93150c9-eb99-4c62-8fc1-c414c8a0f78d', 'drv-003-siti',  NOW(), NOW())
ON CONFLICT (id) DO UPDATE
SET company_id = EXCLUDED.company_id,
    admin_id   = EXCLUDED.admin_id;