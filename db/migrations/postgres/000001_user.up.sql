
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
    user_id varchar(255) NOT NULL,
    organization_id varchar(255),
    phone_number varchar(100) NULL,
    email varchar(510) NULL,
    first_name varchar(1020) NULL,
    last_name varchar(1020) NULL,
    company_type varchar(255) NULL, -- ppip / dkp
    address varchar(1020) NULL,
    domisili varchar(1020) NULL,
    nib varchar(1020) NULL,
    npwp varchar(1020) NULL,
    pic_name varchar(1020) NULL,
    pic_phone varchar(1020) NULL,
    pic_email varchar(1020) NULL, 
    akta_perusahaan varchar(1020) NULL,
    nib_file varchar(1020) NULL,
    tdp varchar(1020) NULL,
    ktp varchar(1020) NULL,
    npwp_perusahaan varchar(1020) NULL,
    surat_kuasa varchar(1020) NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS customer (
    id varchar(255) PRIMARY KEY,
    user_id varchar(255) NOT NULL,
    organization_id varchar(255),
    company_id varchar(255),
    sim_status varchar(255),
    sim_number varchar(255),
    approval_status varchar(255),
    customer_id varchar(255),
    phone_number varchar(100),
    email varchar(510),
    first_name varchar(1020),
    last_name varchar(1020),
    nickname varchar(1020),
    referral_code varchar(255),
    customer_id_parent varchar(255),
    date_start date,
    place_of_birth varchar(255),
    date_of_birth date,
    country_of_birth varchar(255),
    mother_name varchar(255),
    normal_retirement_age varchar(30),
    citizenship varchar(255),
    sex varchar(255),
    marital_status varchar(255),
    occupation varchar(255),
    position varchar(255),
    address varchar(1020),
    mailing_address varchar(1020),
    office_address varchar(1020),
    phone_office varchar(100),
    mobile_phone varchar(100),
    source_of_funds varchar(255),
    annual_income varchar(255),
    purpose_of_opening_account varchar(255),
    name_on_bank_account varchar(255),
    bank_account_number varchar(255),
    bank_name varchar(255),
    identification_number varchar(255),
    tax_identification_number varchar(255),
    employer_percentage varchar(255),
    employer_amount varchar(255),
    customer_percentage varchar(255),
    customer_amount varchar(255),
    effective_date date,
    payment_method varchar(255),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);