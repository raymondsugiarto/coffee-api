ALTER TABLE IF EXISTS approval ADD COLUMN detail varchar(255);
ALTER TABLE IF EXISTS approval ADD COLUMN reason varchar(255);
ALTER TABLE IF EXISTS approval ADD COLUMN created_date date;
