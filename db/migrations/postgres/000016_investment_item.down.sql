DROP TABLE IF EXISTS investment_item;
ALTER TABLE IF EXISTS investment ADD COLUMN investment_type varchar(255) NULL;
ALTER TABLE IF EXISTS investment ADD COLUMN investment_product_id varchar(255) NULL;
ALTER TABLE IF EXISTS investment DROP COLUMN investment_at;

DROP TABLE IF EXISTS investment_distribution;
