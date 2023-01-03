package database

import (
	"context"
	"log"

	"github.com/google/uuid"
)

// BankClient - DTO for client table
type BankClient struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Surname   string    `db:"surname"`
	BirthDate string    `db:"birth_date"`
	Passport  string    `db:"passport"`
	Job       string    `db:"job"`
}

// GetBankClientByID returns bankClient by bank_client_id
func GetBankClientByID(ctx context.Context, dbc SQLExecutor, clientID uuid.UUID) (*BankClient, error) {
	var bankClient *BankClient
	err := dbc.SelectContext(ctx, &bankClient, `
SELECT 
	id, 
	questionary_id, 
	created_at,
	created_by, 
	modified_at,
	modified_by,
	text,
	rating
FROM comments WHERE questionary_id = $1`, clientID)
	if err != nil {
		log.Printf("failed get comments from database")
		return nil, err
	}
	return bankClient, nil
}

// AddNewBankClient - creator for DTO
func AddNewBankClient(ctx context.Context, dbc SQLExecutor, bankClient *BankClient) error {
	return execInsertObjectQuery(ctx, dbc, `
INSERT INTO bank_clients (
	id, name, surname, birth_date, passport, job
) VALUES (
    :id, :name, :surname, :birth_date, :passport, :job
)`, bankClient)
}
