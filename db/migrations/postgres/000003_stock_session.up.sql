-- ============================================================
-- Coffee Cart Driver Inventory Management Module
-- Migration 000003
-- ============================================================

-- Extend existing `item` table with fields needed for driver stock sessions
ALTER TABLE item
    ADD COLUMN IF NOT EXISTS sku varchar(100) NULL,
    ADD COLUMN IF NOT EXISTS category_id varchar(255) NULL,
    ADD COLUMN IF NOT EXISTS cost_price numeric(20, 4) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT TRUE;

CREATE INDEX IF NOT EXISTS idx_item_sku ON item(sku);
CREATE INDEX IF NOT EXISTS idx_item_category ON item(category_id);
CREATE INDEX IF NOT EXISTS idx_item_active ON item(is_active);

-- Item category (replaces product_category)
CREATE TABLE IF NOT EXISTS item_category (
    id varchar(255) PRIMARY KEY,
    organization_id varchar(255) NULL,
    name varchar(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX IF NOT EXISTS idx_item_category_org ON item_category(organization_id);

-- Stock Session: one per (employee, date)
CREATE TABLE IF NOT EXISTS stock_session (
    id varchar(255) PRIMARY KEY,
    organization_id varchar(255) NULL,
    employee_id varchar(255) NOT NULL,
    date DATE NOT NULL,
    status varchar(50) NOT NULL DEFAULT 'OPEN',
    opened_at TIMESTAMP NOT NULL,
    closed_at TIMESTAMP NULL,
    total_sales numeric(20, 4) NOT NULL DEFAULT 0,
    total_cash numeric(20, 4) NOT NULL DEFAULT 0,
    total_qris numeric(20, 4) NOT NULL DEFAULT 0,
    total_other numeric(20, 4) NOT NULL DEFAULT 0,
    total_payment numeric(20, 4) NOT NULL DEFAULT 0,
    difference numeric(20, 4) NOT NULL DEFAULT 0,
    total_items INT NOT NULL DEFAULT 0,
    notes TEXT NULL,
    created_by varchar(255) NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT uq_stock_session_employee_date UNIQUE (employee_id, date)
);

CREATE INDEX IF NOT EXISTS idx_stock_session_org ON stock_session(organization_id);
CREATE INDEX IF NOT EXISTS idx_stock_session_date ON stock_session(date);
CREATE INDEX IF NOT EXISTS idx_stock_session_status ON stock_session(status);

CREATE TABLE IF NOT EXISTS stock_session_item (
    id varchar(255) PRIMARY KEY,
    session_id varchar(255) NOT NULL,
    item_id varchar(255) NOT NULL,
    out_qty INT NOT NULL DEFAULT 0,
    return_qty INT NOT NULL DEFAULT 0,
    sold_qty INT NOT NULL DEFAULT 0,
    selling_price_snapshot numeric(20, 4) NOT NULL DEFAULT 0,
    cost_price_snapshot numeric(20, 4) NOT NULL DEFAULT 0,
    subtotal numeric(20, 4) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_ssi_session FOREIGN KEY (session_id)
        REFERENCES stock_session(id) ON DELETE CASCADE,
    CONSTRAINT fk_ssi_item FOREIGN KEY (item_id)
        REFERENCES item(id) ON DELETE RESTRICT,
    CONSTRAINT chk_ssi_return_qty CHECK (return_qty <= out_qty),
    CONSTRAINT chk_ssi_qty_nonneg CHECK (out_qty >= 0 AND return_qty >= 0 AND sold_qty >= 0)
);

CREATE INDEX IF NOT EXISTS idx_ssi_session ON stock_session_item(session_id);
CREATE INDEX IF NOT EXISTS idx_ssi_item ON stock_session_item(item_id);

CREATE TABLE IF NOT EXISTS payment_detail (
    id varchar(255) PRIMARY KEY,
    session_id varchar(255) NOT NULL,
    payment_method varchar(50) NOT NULL,
    amount numeric(20, 4) NOT NULL DEFAULT 0,
    reference_number varchar(255) NULL,
    notes TEXT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_pd_session FOREIGN KEY (session_id)
        REFERENCES stock_session(id) ON DELETE CASCADE,
    CONSTRAINT chk_pd_amount_nonneg CHECK (amount >= 0)
);

CREATE INDEX IF NOT EXISTS idx_pd_session ON payment_detail(session_id);
CREATE INDEX IF NOT EXISTS idx_pd_method ON payment_detail(payment_method);

CREATE TABLE IF NOT EXISTS cash_adjustment (
    id varchar(255) PRIMARY KEY,
    session_id varchar(255) NOT NULL,
    type varchar(50) NOT NULL,
    amount numeric(20, 4) NOT NULL DEFAULT 0,
    reason TEXT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_ca_session FOREIGN KEY (session_id)
        REFERENCES stock_session(id) ON DELETE CASCADE,
    CONSTRAINT chk_ca_amount_nonneg CHECK (amount >= 0)
);

CREATE INDEX IF NOT EXISTS idx_ca_session ON cash_adjustment(session_id);

CREATE TABLE IF NOT EXISTS session_log (
    id varchar(255) PRIMARY KEY,
    session_id varchar(255) NOT NULL,
    action varchar(100) NOT NULL,
    admin_id varchar(255) NULL,
    detail TEXT NULL,
    created_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_sl_session FOREIGN KEY (session_id)
        REFERENCES stock_session(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_sl_session ON session_log(session_id);
CREATE INDEX IF NOT EXISTS idx_sl_action ON session_log(action);

-- ============================================================
-- SEED DATA (for development/testing)
-- ============================================================

INSERT INTO item_category (id, organization_id, name, created_at)
VALUES
    ('cat-001-coffee', NULL, 'Coffee', NOW()),
    ('cat-002-non-coffee', NULL, 'Non-Coffee', NOW()),
    ('cat-003-snack', NULL, 'Snack', NOW())
ON CONFLICT (id) DO NOTHING;

INSERT INTO item (id, organization_id, category_id, code, sku, name, price, cost_price, is_active, created_at)
VALUES
    ('prod-001-americano', NULL, 'cat-001-coffee', 'AMR-01', 'AMR-01', 'Americano', 18000, 6000, TRUE, NOW()),
    ('prod-002-latte', NULL, 'cat-001-coffee', 'LAT-01', 'LAT-01', 'Latte', 25000, 9000, TRUE, NOW()),
    ('prod-003-cappuccino', NULL, 'cat-001-coffee', 'CAP-01', 'CAP-01', 'Cappuccino', 25000, 9000, TRUE, NOW()),
    ('prod-004-matcha', NULL, 'cat-002-non-coffee', 'MTC-01', 'MTC-01', 'Matcha Latte', 28000, 11000, TRUE, NOW()),
    ('prod-005-chocolate', NULL, 'cat-002-non-coffee', 'CHC-01', 'CHC-01', 'Chocolate', 22000, 8000, TRUE, NOW()),
    ('prod-006-teh', NULL, 'cat-002-non-coffee', 'TEH-01', 'TEH-01', 'Milk Tea', 20000, 7000, TRUE, NOW()),
    ('prod-007-croissant', NULL, 'cat-003-snack', 'CRS-01', 'CRS-01', 'Croissant', 15000, 5000, TRUE, NOW()),
    ('prod-008-donut', NULL, 'cat-003-snack', 'DNT-01', 'DNT-01', 'Donut', 12000, 4000, TRUE, NOW())
ON CONFLICT (id) DO NOTHING;

-- Seed employees (drivers) if not exist — uses existing admin table
INSERT INTO admin (id, organization_id, user_id, admin_type, phone_number, email, first_name, last_name, created_at)
VALUES
    ('drv-001-budi', NULL, NULL, 'EMPLOYEE', '0812000001', 'budi@coffee.test', 'Budi', 'Santoso', NOW()),
    ('drv-002-andi', NULL, NULL, 'EMPLOYEE', '0812000002', 'andi@coffee.test', 'Andi', 'Wijaya', NOW()),
    ('drv-003-siti', NULL, NULL, 'EMPLOYEE', '0812000003', 'siti@coffee.test', 'Siti', 'Nurhaliza', NOW())
ON CONFLICT (id) DO NOTHING;
