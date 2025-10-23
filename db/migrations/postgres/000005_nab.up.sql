-- NAB = Nilai Aktiva Bersih
CREATE TABLE IF NOT EXISTS net_asset_value (
    id varchar(255) PRIMARY KEY,
    organization_id varchar(255),
    investment_product_id varchar(255) NULL,
    created_date date,
    amount numeric(20,4) NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS portfolio (
    id varchar(255) PRIMARY KEY,
    organization_id varchar(255),
    company_id varchar(255) NULL,
    customer_id varchar(255) NULL,
    investment_product_id varchar(255) NULL,
    portfolio_date date,
    amount numeric(20,4) NULL,
    net_asset_value_amount numeric(20,4) NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);
