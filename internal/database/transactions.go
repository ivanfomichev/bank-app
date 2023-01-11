package database

import (
	"context"
	"log"

	"github.com/google/uuid"
)

// Transaction - DTO for transaction table
type Transaction struct {
	TransactionID uuid.UUID `json:"transaction_id"`
	AccountID     uuid.UUID `json:"account_id" validate:"required"`
	AccountToID   uuid.UUID `json:"account_to_id"`
	Amount        int32     `json:"amount" validate:"required"`
	TrType        string    `json:"tr_type" validate:"required"`
}

// GetTransactions returns transactions
func GetTransactions(ctx context.Context, dbc SQLExecutor) ([]*Transaction, error) {
	transactions := make([]*Transaction, 0)
	rows, err := dbc.QueryContext(
		ctx,
		`SELECT * FROM transactions`,
	)
	if err != nil {
		log.Printf("failed get transactions from database")
		return nil, err
	}
	for rows.Next() {
		var transaction_id uuid.UUID
		var account_id uuid.UUID
		var account_to_id uuid.UUID
		var amount int32
		var tr_type string
		err := rows.Scan(&transaction_id, &account_id, &account_to_id, &amount, &tr_type)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		transactions = append(transactions, &Transaction{
			TransactionID: transaction_id,
			AccountID:     account_id,
			AccountToID:   account_to_id,
			Amount:        amount,
			TrType:        tr_type,
		})
	}
	return transactions, nil
}

// AddNewTransaction - creator for DTO
func AddNewTransaction(ctx context.Context, dbc SQLExecutor, transaction *Transaction) error {
	_, err := dbc.ExecContext(ctx, `
INSERT INTO transactions (
	transaction_id, account_id, account_to_id, amount, tr_type
) VALUES ($1, $2, $3, $4, $5)`,
		transaction.TransactionID,
		transaction.AccountID,
		transaction.AccountToID,
		transaction.Amount,
		transaction.TrType,
	)
	return err
}
