CREATE TABLE IF NOT EXISTS user_identity_verification (
    id varchar(255) PRIMARY KEY,
    organization_id varchar(255),
    user_id varchar(100) NULL,
    identity_for varchar(255) NULL,
    identity_type varchar(255) NULL,
    user_identity varchar(255) NULL,
    unique_code varchar(255) NULL,
    try_count varchar(255) NULL,
    status varchar(255) NULL,
    expired_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);
