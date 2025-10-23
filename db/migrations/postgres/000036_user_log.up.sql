CREATE TABLE IF NOT EXISTS "user_log" (
    id varchar(255) PRIMARY KEY,
    organization_id varchar(255),
    company_id varchar(255),
    user_credential_id varchar(255),
    ref_module varchar(255),
    ref_table varchar(255),
    ref_id varchar(255),
    description text,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);