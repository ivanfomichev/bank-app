package database

import (
	"context"
	"log"

	"github.com/google/uuid"
)

// BankClient - DTO for client table
type BankClient struct {
	ID uuid.UUID `db:"id"`
}

// GetBankClientByID returns bankClient by bank_client_id
func GetBankClientByID(ctx context.Context, dbc SQLExecutor, clientID string) (*BankClient, error) {
	bankClient := new(BankClient)
	err := dbc.QueryRowxContext(ctx,
		`SELECT * FROM bank_clients WHERE id = $1`,
		clientID,
	).StructScan(bankClient)
	if err != nil {
		log.Printf("failed get client from database")
		return nil, err
	}
	return bankClient, nil
}

// AddNewBankClient - creator for DTO
func AddNewBankClient(ctx context.Context, dbc SQLExecutor, bankClient *BankClient) error {
	return execInsertObjectQuery(ctx,
		dbc,
		`INSERT INTO bank_clients (id) VALUES (:id)`,
		bankClient,
	)
}
