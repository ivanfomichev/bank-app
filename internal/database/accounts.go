package database

import (
	"context"
	"log"

	"github.com/google/uuid"
)

// Account - DTO for account table
type Account struct {
	AccountID uuid.UUID `db:"account_id"`
	ClientID  uuid.UUID `db:"client_id"`
	Currency  string    `db:"currency" validate:"required"`
	Balance   int32     `db:"balance" validate:"required"`
}

// GetAccount returns account by account_id
func GetAccountByID(ctx context.Context, dbc SQLExecutor, accountID string) (*Account, error) {
	account := new(Account)
	err := dbc.QueryRowxContext(ctx,
		`SELECT * FROM accounts WHERE account_id = $1`,
		accountID,
	).StructScan(account)
	if err != nil {
		log.Printf("failed get account from database")
		return nil, err
	}
	return account, nil
}

// AddNewBankClient - creator for DTO
func AddNewAccount(ctx context.Context, dbc SQLExecutor, account *Account) error {
	err := execInsertObjectQuery(ctx,
		dbc,
		`INSERT INTO accounts (
			account_id, client_id, currency, balance
		) VALUES (
			:account_id, :client_id, :currency, :balance
		)`,
		account,
	)
	return err
}

// UpdateAccountByID - updater for DTO
func UpdateAccountByID(ctx context.Context, dbc SQLExecutor, clientID uuid.UUID, balance int32) error {
	return updateTableColWithProvidedKey(ctx,
		dbc,
		`UPDATE accounts SET balance = $1 WHERE account_id = $2`,
		balance,
		clientID,
	)
}
