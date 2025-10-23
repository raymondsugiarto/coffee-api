ALTER TABLE IF EXISTS pension_benefit_recipient ADD COLUMN phone_number varchar(255);

ALTER TABLE IF EXISTS customer DROP COLUMN annual_income;
ALTER TABLE IF EXISTS customer ADD COLUMN annual_income varchar(255);
