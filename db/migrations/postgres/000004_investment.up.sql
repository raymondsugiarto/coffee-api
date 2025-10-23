CREATE TABLE IF NOT EXISTS investment_product (
    id varchar(255) PRIMARY KEY,
    organization_id varchar(255),
    code varchar(255) NULL,
    name varchar(1020) NULL,
    description text NULL,
    fix_income numeric(20,4),
    stock_value numeric(20,4),
    mixed_value numeric(20,4),
    money_market numeric(20,4),
    sharia_value numeric(20,4),
    fund_fact_sheet text NULL,
    riplay text NULL,
    admin_fee numeric(20,4),
    management_fee numeric(20,4),
    founder_fee numeric(20,4),
    commission_fee numeric(20,4),
    status varchar(255) NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS investment (
    id varchar(255) PRIMARY KEY,
    organization_id varchar(255),
    company_id varchar(100) NULL,
    customer_id varchar(255) NULL,
    investment_type varchar(255) NULL,
    investment_product_id varchar(255) NULL,
    amount numeric(20,4) NULL,
    status varchar(255) NULL,
    expired_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);
