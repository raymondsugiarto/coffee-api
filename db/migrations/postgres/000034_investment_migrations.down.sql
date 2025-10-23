ALTER TABLE benefit_participation DROP COLUMN investment_id;
ALTER TABLE investment DROP COLUMN source;
ALTER TABLE investment_distribution 
    DROP COLUMN base_contribution,
    DROP COLUMN voluntary_contribution;

