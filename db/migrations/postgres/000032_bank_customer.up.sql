CREATE TABLE IF NOT EXISTS bank_customer (
    id varchar(255) PRIMARY KEY,
    organization_id varchar(255),
    customer_id varchar(255),
    bank_code varchar(255) NULL,
    bank_name varchar(255) NULL,
    account_name varchar(255) NULL,
    account_number varchar(255) NULL,
    is_default boolean NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);