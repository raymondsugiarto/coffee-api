-- Reverse 000014 accounting. account_mutation has to drop first
-- because of its FK to account(id).

DROP INDEX IF EXISTS idx_account_mutation_ref;
DROP INDEX IF EXISTS idx_account_mutation_org_date;
DROP INDEX IF EXISTS idx_account_mutation_account_date;
DROP TABLE IF EXISTS account_mutation;

DROP INDEX IF EXISTS idx_account_name;
DROP INDEX IF EXISTS idx_account_org;
DROP TABLE IF EXISTS account;
