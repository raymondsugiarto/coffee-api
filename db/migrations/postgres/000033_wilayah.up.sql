CREATE TABLE country (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    cca2 VARCHAR(5) NOT NULL,
    cca3 VARCHAR(5) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE province (
    id VARCHAR(255) PRIMARY KEY,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE regency (
    id VARCHAR(255) PRIMARY KEY,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    province_id VARCHAR(255) REFERENCES province(id),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE district (
    id VARCHAR(255) PRIMARY KEY,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    regency_id VARCHAR(255) REFERENCES regency(id),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE village (
    id VARCHAR(255) PRIMARY KEY,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    postal_code VARCHAR(50),
    district_id VARCHAR(255) REFERENCES district(id),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

-- Add unique constraints to ensure no duplicate codes
ALTER TABLE province
ADD CONSTRAINT unique_province_code UNIQUE (code);

ALTER TABLE regency
ADD CONSTRAINT unique_regency_code UNIQUE (code);

ALTER TABLE district
ADD CONSTRAINT unique_district_code UNIQUE (code);

ALTER TABLE village
ADD CONSTRAINT unique_village_code UNIQUE (code);