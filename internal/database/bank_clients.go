package database

import (
	"context"
	"log"

	"github.com/google/uuid"
)

// BankClient - DTO for client table
type BankClient struct {
	ClientID      uuid.UUID `db:"client_id"`
	IdentityField int32     `db:"identity_field"`
	ClientName    string    `db:"client_name"`
}

// GetBankClientByID returns bankClient by client_id
func GetBankClientByID(ctx context.Context, dbc SQLExecutor, clientID string) (*BankClient, error) {
	bankClient := new(BankClient)
	err := dbc.QueryRowxContext(ctx,
		`SELECT * FROM bank_clients WHERE client_id = $1`,
		clientID,
	).StructScan(bankClient)
	if err != nil {
		log.Printf("failed to get bank client from database")
		return nil, err
	}
	return bankClient, nil
}

// GetBankClientByIdentity returns bankClient by identity_field
func GetBankClientByIdentity(ctx context.Context, dbc SQLExecutor, identityField int32) error {
	bankClient := new(BankClient)
	row := dbc.QueryRowxContext(ctx,
		`SELECT * FROM bank_clients WHERE identity_field = $1`,
		identityField,
	)
	err := row.Scan(bankClient)
	if err != nil {
		log.Printf("failed to get bank client from database")
		return err
	}
	return nil
}

// AddNewBankClient - creator for DTO
func AddNewBankClient(ctx context.Context, dbc SQLExecutor, bankClient *BankClient) error {
	return execInsertObjectQuery(ctx,
		dbc,
		`INSERT INTO bank_clients (
			identity_field, client_name
		) VALUES (
			:identity_field, :client_name
		)`,
		bankClient,
	)
}
