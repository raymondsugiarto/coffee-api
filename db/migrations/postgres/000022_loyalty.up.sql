CREATE TABLE IF NOT EXISTS customer_point (
    id varchar(255) PRIMARY KEY,
    organization_id varchar(255),
    customer_id varchar(100) NULL,
    point numeric(20,4) NULL,
    direction varchar(255) NULL,
    description varchar(255) NULL,
    ref_id varchar(255) NULL,
    ref_code varchar(255) NULL,
    ref_module varchar(255) NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);
