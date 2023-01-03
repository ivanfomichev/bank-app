-- +migrate Up
CREATE TABLE bank_clients
(
    id         uuid PRIMARY KEY,
    name       VARCHAR(256) NOT NULL,
    surname    VARCHAR(256) NOT NULL,
    birth_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    passport   VARCHAR(256) NOT NULL,
    job        VARCHAR(256)
);

CREATE TYPE currencies AS ENUM (
    'USD',
    'COP',
    'MXN'
    );

CREATE TABLE accounts
(
    id             uuid PRIMARY KEY,
    bank_client_id uuid NOT NULL,
    currency       currencies NOT NULL,
    balance        VARCHAR(256) NOT NULL
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
    id         uuid PRIMARY KEY,
    account_id uuid NOT NULL,
    tr_date    TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    tr_type    transaction_type NOT NULL,
    tr_status  transaction_status NOT NULL

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
-- +migrate StatementEnd
