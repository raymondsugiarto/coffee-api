CREATE TABLE IF NOT EXISTS item (
    id varchar(255) PRIMARY KEY,
    organization_id varchar(255) NULL,
    code varchar(100) NULL,
    name varchar(255) NULL,
    price numeric(20, 4) NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS item_company (
    id varchar(255) PRIMARY KEY,
    organization_id varchar(255) NULL,
    item_id varchar(255) NULL,
    company_id varchar(255) NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS "order" (
    id varchar(255) PRIMARY KEY,
    organization_id varchar(255) NULL,
    company_id varchar(255) NULL,
    admin_id varchar(255) NULL,
    customer_id varchar(255) NULL,
    code varchar(100) NULL,
    order_at TIMESTAMP NOT NULL,
    total_qty INT NOT NULL,
    total_amount numeric(20, 4) NULL,
    status varchar(100) NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);


CREATE TABLE IF NOT EXISTS order_item (
    id varchar(255) PRIMARY KEY,
    order_id varchar(255) NULL,
    item_id varchar(255) NULL,
    qty INT NOT NULL,
    price numeric(20, 4) NULL,
    subtotal numeric(20, 4) NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);


CREATE TABLE IF NOT EXISTS order_payment (
    id varchar(255) PRIMARY KEY,
    order_id varchar(255) NULL,
    payment_method_code varchar(255) NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);


