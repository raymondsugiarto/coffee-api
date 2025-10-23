ALTER TABLE IF EXISTS portfolio ADD COLUMN participant_id VARCHAR(255);
ALTER TABLE IF EXISTS portfolio ADD COLUMN ip numeric(20, 4);

ALTER TABLE IF EXISTS portfolio DROP COLUMN amount;
ALTER TABLE IF EXISTS portfolio DROP COLUMN net_asset_value_amount;

ALTER TABLE IF EXISTS unit_link ADD COLUMN transaction_date date;
ALTER TABLE IF EXISTS unit_link ADD COLUMN organization_id VARCHAR(255);
ALTER TABLE IF EXISTS unit_link ALTER COLUMN ip TYPE numeric(20, 4) USING ip::numeric(20, 4);