ALTER TABLE IF EXISTS investment_item ADD COLUMN employer_amount numeric(20,4) DEFAULT 0.0;
ALTER TABLE IF EXISTS investment_item ADD COLUMN employee_amount numeric(20,4) DEFAULT 0.0;
ALTER TABLE IF EXISTS investment_item ADD COLUMN voluntary_amount numeric(20,4) DEFAULT 0.0;
ALTER TABLE IF EXISTS investment_item ADD COLUMN education_fund_amount numeric(20,4) DEFAULT 0.0;