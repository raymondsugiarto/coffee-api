CREATE TABLE IF NOT EXISTS ticket (
    id VARCHAR(255) PRIMARY KEY,
    organization_id VARCHAR(255) NULL,
    user_id VARCHAR(255) NULL,
    title VARCHAR(1020) NULL,
    message text NULL,
    status VARCHAR(255) NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE
);