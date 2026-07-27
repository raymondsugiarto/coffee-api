-- ============================================================
-- 000014: accounting module
-- ============================================================
-- Two tables:
--
--   account          : master chart-of-accounts (operator-managed
--                      categories like "Kas", "Piutang", "Penjualan",
--                      "Beban Gaji", etc.). CRUD.
--
--   account_mutation : append-only ledger of movements tied to a
--                      reference row elsewhere in the system.
--                      Examples:
--                        ref_table = 'stock_session'
--                        ref_id    = <stock_session.id>
--                        account_id = account "Penjualan"
--                        amount    = total_payment
--                      This table is the single source of truth for
--                      period totals so reports can run as SELECT
--                      SUM(amount) GROUP BY account_id, date.
--
-- We deliberately do NOT wire foreign keys from
-- account_mutation.ref_id to the upstream tables (orders,
-- stock_sessions, etc.). The upstream schema is owned by other
-- modules and a hard FK would couple accounting to their migration
-- timeline. ref_table + ref_id are a soft pointer; integrity is
-- enforced at write time by the service layer that produces the
-- mutation.

CREATE TABLE IF NOT EXISTS account (
    id              varchar(255) PRIMARY KEY,
    organization_id varchar(255),
    name            varchar(255) NOT NULL,
    code            varchar(64)  NOT NULL,
    created_at      TIMESTAMP    NOT NULL,
    updated_at      TIMESTAMP    NULL,
    deleted_at      TIMESTAMP    NULL
);

CREATE TABLE IF NOT EXISTS account_mutation (
    id              varchar(255)   PRIMARY KEY,
    organization_id varchar(255),
    account_id      varchar(255)   NOT NULL,
    -- Amount is signed: positive = debit (asset/expense inflow
--                    into this account), negative = credit
--                    (asset/expense outflow). Reports interpret
--                    by account type at query time.
    amount          numeric(20, 4) NOT NULL,
    description     text           NULL,
    ref_id          varchar(255)   NOT NULL,
    ref_table       varchar(64)    NOT NULL,
    ref_module      varchar(64)    NOT NULL,
    created_at      TIMESTAMP      NOT NULL,
    updated_at      TIMESTAMP      NULL,
    deleted_at      TIMESTAMP      NULL
);
