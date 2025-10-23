ALTER TABLE IF EXISTS company ADD COLUMN company_code varchar(10);
ALTER TABLE IF EXISTS company ADD COLUMN agreement_fee numeric(20,4);
ALTER TABLE IF EXISTS company ADD COLUMN cooperation_agreement varchar(1020);

ALTER TABLE IF EXISTS customer DROP COLUMN annual_income;
ALTER TABLE IF EXISTS customer ADD COLUMN annual_income numeric(20,4);

ALTER TABLE IF EXISTS customer DROP COLUMN employer_percentage;
ALTER TABLE IF EXISTS customer ADD COLUMN employer_percentage numeric(20,4);

ALTER TABLE IF EXISTS customer DROP COLUMN employer_amount;
ALTER TABLE IF EXISTS customer ADD COLUMN employer_amount numeric(20,4);

ALTER TABLE IF EXISTS customer DROP COLUMN customer_percentage;
ALTER TABLE IF EXISTS customer ADD COLUMN customer_percentage numeric(20,4);

ALTER TABLE IF EXISTS customer DROP COLUMN customer_amount;
ALTER TABLE IF EXISTS customer ADD COLUMN customer_amount numeric(20,4);