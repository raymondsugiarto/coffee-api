CREATE TABLE IF NOT EXISTS role (
    id varchar(255) PRIMARY KEY,
    organization_id varchar(255),
    code varchar(100) NULL,
    name varchar(255) NULL,
    role_id_parent varchar(255) NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);


CREATE TABLE IF NOT EXISTS user_has_role (
    id varchar(255) PRIMARY KEY,
    user_id varchar(255),
    role_id varchar(255),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);


CREATE TABLE IF NOT EXISTS permission (
    id varchar(255) PRIMARY KEY,
    organization_id varchar(255),
    code varchar(255) NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS role_permission (
    id varchar(255) PRIMARY KEY,
    role_id varchar(255),
    permission_id varchar(255),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);