CREATE TABLE IF NOT EXISTS benefit_type (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    minimum_time_period_months INT NOT NULL,
    minimum_contribution NUMERIC(10, 2) NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS benefit_participation (
    id VARCHAR(255) PRIMARY KEY,
    organization_id VARCHAR(255) NOT NULL,
    customer_id VARCHAR(255) NOT NULL,
    participant_id VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL,

    external_dplk_name VARCHAR(255),
    external_dplk_participant_number VARCHAR(255),
    external_dplk_monthly_contribution NUMERIC(10, 2),
    has_bpjs_pension_program BOOLEAN,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS benefit_participation_detail (
    id VARCHAR(255) PRIMARY KEY,
    benefit_participation_id VARCHAR(255) NOT NULL,
    benefit_type_id VARCHAR(255) NOT NULL,
    time_period_months INT NOT NULL,
    planned_withdrawal_months INT NOT NULL,
    monthly_contribution NUMERIC(10, 2) NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE
);