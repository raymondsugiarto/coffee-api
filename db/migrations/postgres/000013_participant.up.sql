CREATE TABLE IF NOT EXISTS participant (
    id varchar(255) PRIMARY KEY,
    organization_id varchar(255),
    code varchar(255),
    customer_id varchar(255),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

ALTER TABLE investment ADD COLUMN participant_id varchar(255);