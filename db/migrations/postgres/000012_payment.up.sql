CREATE TABLE IF NOT EXISTS bank (
    id varchar(255) PRIMARY KEY,
    organization_id varchar(255),
    bank_code varchar(255) NULL,
    bank_name varchar(255) NULL,
    account_name varchar(255) NULL,
    account_number varchar(255) NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS investment_payment (
    id varchar(255) PRIMARY KEY,
    organization_id varchar(255),
    investment_id varchar(255),
    payment_method varchar(255) NULL,
    status varchar(255) NULL,
    amount decimal(20, 2) NULL,
    bank_id varchar(255) NULL,
    bank_code varchar(255) NULL,
    bank_name varchar(255) NULL,
    account_name varchar(255) NULL,
    account_number varchar(255) NULL,
    confirmation_image_url varchar(255) NULL,
    payment_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);