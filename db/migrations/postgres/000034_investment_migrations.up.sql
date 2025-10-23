ALTER TABLE investment ADD COLUMN source VARCHAR(255);
ALTER TABLE benefit_participation ADD COLUMN investment_id VARCHAR(36);
ALTER TABLE investment_distribution 
    ADD COLUMN base_contribution NUMERIC(15,2),
    ADD COLUMN voluntary_contribution NUMERIC(15,2);

