CREATE TABLE IF NOT EXISTS notification (
    id VARCHAR(255) PRIMARY KEY,
    organization_id VARCHAR(255) NULL,
    user_id VARCHAR(255) NULL,
    ref_module VARCHAR(255),
    ref_table VARCHAR(255),
    ref_id VARCHAR(255) ,
    ref_code VARCHAR(255) ,
    description text,
    notify_at TIMESTAMP WITH TIME ZONE,
    read_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE IF EXISTS investment ADD COLUMN code VARCHAR(255);