CREATE TABLE IF NOT EXISTS pension_benefit_recipient (
    id varchar(255) PRIMARY KEY,
    customer_id varchar(255) NULL,
	name varchar(255) NULL,
	relationship varchar(255) NULL,
	date_of_birth varchar(255) NULL,
	country_of_birth varchar(255) NULL,
	identification_number varchar(255) NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);
