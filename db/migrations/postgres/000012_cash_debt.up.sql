-- ============================================================
-- 000012: cash_debt table
-- ============================================================
-- Stand-alone cash debt ledger. Each row is one driver-issued
-- cash advance tied to an employee, a date, a nominal amount,
-- and a payment method (CASH = physical cash handed to the
-- driver; CASHLESS = settled via QRIS/TRANSFER/OTHER so the
-- cash float is not depleted).
--
-- Schema mirrors the existing admin/payment conventions: id,
-- organization scope, soft-delete via deleted_at, plus
-- created_at / updated_at. Indexes on employee + date so the
-- most common listing / filter paths hit a single index.

CREATE TABLE IF NOT EXISTS cash_debt (
    id              varchar(255) PRIMARY KEY,
    organization_id varchar(255),
    admin_id_employee varchar(255) NOT NULL,
    date            date           NOT NULL,
    amount          numeric(20, 4) NOT NULL DEFAULT 0,
    payment_method  varchar(20)    NOT NULL DEFAULT 'CASH',
    notes           text           NULL,
    created_at      TIMESTAMP      NOT NULL,
    updated_at      TIMESTAMP      NULL,
    deleted_at      TIMESTAMP      NULL
);

