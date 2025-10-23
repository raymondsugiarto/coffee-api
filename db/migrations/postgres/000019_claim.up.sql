CREATE TABLE claim (
    id VARCHAR(36) PRIMARY KEY,
    organization_id VARCHAR(36),
    participant_id VARCHAR(36),
    bank_name VARCHAR(255),
    account_name VARCHAR(255),
    account_number VARCHAR(255),
    bank_branch VARCHAR(255),
    amount NUMERIC(15, 2),
    approval_status VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE
);

