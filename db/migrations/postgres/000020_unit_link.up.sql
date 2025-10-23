CREATE TABLE unit_link (
    id VARCHAR(36) PRIMARY KEY,
    participant_id VARCHAR(36) NOT NULL,
    customer_id VARCHAR(36),
    investment_product_id VARCHAR(36) NOT NULL,
    type VARCHAR(255) NOT NULL,
    total_amount NUMERIC(18, 2) NOT NULL,
    nab NUMERIC(10, 6) NOT NULL,
    ip VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);