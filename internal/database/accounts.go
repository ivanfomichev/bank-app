package database

import (
	"context"
	"log"

	"github.com/google/uuid"
)

// Account - DTO for account table
type Account struct {
	ID           uuid.UUID `db:"id"`
	BankClientID string    `db:"bank_client_id"`
	Currency     string    `db:"currency"`
	Balance      *int32    `db:"balance"`
}

// GetAccount returns account by account_id
func GetAccountByID(ctx context.Context, dbc SQLExecutor, accountID string) (*Account, error) {
	account := new(Account)
	err := dbc.QueryRowxContext(ctx,
		`SELECT * FROM accounts WHERE id = $1`,
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
		`INSERT INTO accounts (id, bank_client_id, currency, balance) VALUES (:id, :bank_client_id, :currency, :balance)`,
		account,
	)
	return err
}
