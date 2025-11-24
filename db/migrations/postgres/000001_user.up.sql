CREATE TABLE IF NOT EXISTS organization (
    id varchar(255) PRIMARY KEY,
    code varchar(100) NULL,
    name varchar(255) NULL,
    origin varchar(255) NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS "user" (
    id varchar(255) PRIMARY KEY,
    organization_id varchar(255),
    user_type varchar(255) NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS "admin" (
    id varchar(255) PRIMARY KEY,
    organization_id varchar(255),
    user_id varchar(255) NULL,
    admin_type varchar(100) NULL,
    phone_number varchar(100) NULL,
    email varchar(510) NULL,
    first_name varchar(1020) NULL,
    last_name varchar(1020) NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS "user_credential" (
    id varchar(255) PRIMARY KEY,
    organization_id varchar(255),
    user_id varchar(100) NULL,
    username varchar(510) NULL,
    password varchar(1020) NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS company (
    id varchar(255) PRIMARY KEY,
    organization_id varchar(255),
    phone_number varchar(100) NULL,
    name varchar(300) NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS admin_company (
    id varchar(255) PRIMARY KEY,
    organization_id varchar(255),
    company_id varchar(255) NOT NULL,
    admin_id varchar(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS customer (
    id varchar(255) PRIMARY KEY,
    user_id varchar(255) NOT NULL,
    organization_id varchar(255),
    company_id varchar(255),
    phone_number varchar(100),
    email varchar(510),
    first_name varchar(1020),
    last_name varchar(1020),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);