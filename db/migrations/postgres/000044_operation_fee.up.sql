CREATE TABLE IF NOT EXISTS "transaction_fee" (
    id varchar(255) PRIMARY KEY,
    organization_id varchar(255),
    company_id varchar(255),
    investment_product_id varchar(255),
    participant_id varchar(255),
    type varchar(255),
    transaction_date date ,
    ip numeric(20, 4),
    operation_fee numeric(20, 4),
    nav numeric(20, 4),
    portfolio_amount numeric(20, 4),

    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);