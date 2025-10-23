CREATE TABLE IF NOT EXISTS approval (
    id varchar(255) PRIMARY KEY,
    organization_id varchar(255),
    type varchar(255) NULL,
    action varchar(1020) NULL,
    user_id_request varchar(255) NULL,
    ref_id varchar(255) NULL,
    ref_table varchar(255) NULL,
    status varchar(255) NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);