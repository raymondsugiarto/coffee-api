CREATE TABLE IF NOT EXISTS fee_setting (
    id varchar(255) PRIMARY KEY,
    admin_fee NUMERIC(10, 2) NOT NULL,
    operational_fee NUMERIC(10, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);