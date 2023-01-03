package database

import (
	"context"
	"log"

	"github.com/google/uuid"
)

// Transaction - DTO for transaction table
type Transaction struct {
	ID        uuid.UUID `json:"id"`
	AccountID uuid.UUID `json:"account_id"`
	TrType    string    `json:"tr_type"`
	TrStatus  string    `json:"tr_status"`
}

// GetTransactions returns bankClient by bank_client_id
func GetTransactions(ctx context.Context, dbc SQLExecutor) ([]*Transaction, error) {
	transactions := make([]*Transaction, 0)
	err := dbc.SelectContext(ctx,
		transactions,
		`SELECT * FROM bank_clients WHERE id = $1`,
	)
	if err != nil {
		log.Printf("failed get client from database")
		return nil, err
	}
	return transactions, nil
}

// AddNewTransaction - creator for DTO
func AddNewTransaction(ctx context.Context, dbc SQLExecutor, transaction *Transaction) error {
	return execInsertObjectQuery(ctx,
		dbc,
		`INSERT INTO transactions (id, account_id, tr_type, tr_status) VALUES (:id, :account_id, :tr_type, :tr_status)`,
		transaction,
	)
}
