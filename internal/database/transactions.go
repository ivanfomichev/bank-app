package database

import (
	"context"
	"log"

	"github.com/google/uuid"
)

// Transaction - DTO for transaction table
type Transaction struct {
	ID          uuid.UUID `json:"id"`
	AccountID   uuid.UUID `json:"account_id"`
	AccountToID uuid.UUID `json:"account_to_id"`
	Amount      int32     `json:"amount"`
	TrType      string    `json:"tr_type"`
	TrStatus    string    `json:"tr_status"`
}

// GetTransactions returns transactions
func GetTransactions(ctx context.Context, dbc SQLExecutor) ([]*Transaction, error) {
	transactions := make([]*Transaction, 0)
	rows, err := dbc.QueryContext(ctx,
		`SELECT * FROM transactions`,
		&transactions,
	)
	if err != nil {
		log.Printf("failed get transactions from database")
		return nil, err
	}
	for rows.Next() {
		transaction := Transaction{}
		err := rows.Scan(&transaction)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}
	return transactions, nil
}

// AddNewTransaction - creator for DTO
func AddNewTransaction(ctx context.Context, dbc SQLExecutor, transaction *Transaction) error {
	return execInsertObjectQuery(ctx,
		dbc,
		`INSERT INTO transactions (
			id, account_id, account_to_id, amount, tr_type, tr_status
		) VALUES (
			:id, :account_id, :account_to_id, :amount, :tr_type, :tr_status
		)`,
		transaction,
	)
}
