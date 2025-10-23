ALTER TABLE IF EXISTS company DROP COLUMN company_code;
ALTER TABLE IF EXISTS company DROP COLUMN agreement_fee;
ALTER TABLE IF EXISTS company DROP COLUMN cooperation_agreement;

ALTER TABLE IF EXISTS customer DROP COLUMN annual_income;
ALTER TABLE IF EXISTS customer DROP COLUMN employer_percentage;
ALTER TABLE IF EXISTS customer DROP COLUMN employer_amount;
ALTER TABLE IF EXISTS customer DROP COLUMN customer_percentage;
ALTER TABLE IF EXISTS customer DROP COLUMN customer_amount;