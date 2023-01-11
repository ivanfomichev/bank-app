-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE bank_clients
(
    client_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    identity_field INT CHECK(identity_field > 0),
    client_name VARCHAR(255) NOT NULL,
    UNIQUE (client_id, identity_field)
);

CREATE TYPE currencies AS ENUM (
    'USD',
    'COP',
    'MXN'
);

CREATE TABLE accounts
(
    account_id uuid DEFAULT uuid_generate_v4(),
    client_id  uuid,
    currency   currencies NOT NULL,
    balance    INTEGER NOT NULL CHECK (balance > 0),
    PRIMARY KEY(account_id),
    CONSTRAINT fk_client
        FOREIGN KEY(client_id)
            REFERENCES bank_clients(client_id)
            ON DELETE CASCADE
);
CREATE INDEX accounts_client_id_idx ON accounts (client_id);

CREATE TYPE transaction_type AS ENUM (
    'deposit',
    'withdraw',
    'transfer'
);

CREATE TABLE transactions
(
    transaction_id  uuid DEFAULT uuid_generate_v4(),
    account_id      uuid NOT NULL,
    account_to_id   uuid,
    amount          INTEGER NOT NULL,
    tr_type         transaction_type NOT NULL,
    PRIMARY KEY(transaction_id),
    CONSTRAINT fk_account
        FOREIGN KEY(account_id)
            REFERENCES accounts(account_id)
            ON DELETE CASCADE
);
CREATE INDEX transactions_account_id_idx ON transactions (account_id);
-- +migrate StatementEnd

-- +migrate Down
DROP TABLE bank_clients;
DROP TYPE currencies;
DROP TYPE transaction_type;
DROP TABLE accounts;
DROP TABLE transactions;

DROP EXTENSION IF EXISTS "uuid-ossp";
-- +migrate StatementEnd
