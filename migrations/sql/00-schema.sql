-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE bank_clients
(
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4()
);

CREATE TYPE currencies AS ENUM (
    'USD',
    'COP',
    'MXN'
);

CREATE TABLE accounts
(
    id             uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    bank_client_id uuid NOT NULL,
    currency       currencies NOT NULL,
    balance        INTEGER NOT NULL
);
CREATE INDEX accounts_bank_client_id_idx ON accounts (bank_client_id);

CREATE TYPE transaction_type AS ENUM (
    'deposit',
    'withdraw',
    'transfer'
);

CREATE TYPE transaction_status AS ENUM (
    'done',
    'failed'
);

CREATE TABLE transactions
(
    id              uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id      uuid NOT NULL,
    account_to_id   uuid,
    amount          INTEGER NOT NULL,
    tr_type         transaction_type NOT NULL,
    tr_status       transaction_status NOT NULL
);
CREATE INDEX transactions_account_id_idx ON transactions (account_id);
-- +migrate StatementEnd

-- +migrate Down
DROP TABLE bank_clients;
DROP TYPE currencies;
DROP TYPE transaction_type;
DROP TYPE transaction_status;
DROP TABLE accounts;
DROP TABLE transactions;

DROP EXTENSION IF EXISTS "uuid-ossp";
-- +migrate StatementEnd
