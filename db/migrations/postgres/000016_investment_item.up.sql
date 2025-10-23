CREATE TABLE IF NOT EXISTS investment_item (
    id varchar(255) PRIMARY KEY,
    organization_id varchar(255),
    investment_id varchar(100) NULL,
    customer_id varchar(255) NULL,
    participant_id varchar(255) NULL,
    investment_type varchar(255) NULL,
    investment_product_id varchar(255) NULL,
    type varchar(255) NULL,
    percent numeric(20,4) NULL,
    amount numeric(20,4) NULL,
    status varchar(255) NULL,
    expired_at TIMESTAMP NOT NULL,
    investment_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

ALTER TABLE IF EXISTS investment DROP COLUMN investment_type;
ALTER TABLE IF EXISTS investment DROP COLUMN investment_product_id;

ALTER TABLE IF EXISTS investment ADD COLUMN investment_at TIMESTAMP;
ALTER TABLE IF EXISTS investment ADD COLUMN type varchar(255) NULL;

CREATE TABLE IF NOT EXISTS investment_distribution (
    id varchar(255) PRIMARY KEY,
    organization_id varchar(255),
    type varchar(255) NULL,
    company_id varchar(255) NULL,
    customer_id varchar(255) NULL,
    participant_id varchar(255) NULL,
    investment_product_id varchar(255) NULL,
    percent numeric(20,4) NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);
