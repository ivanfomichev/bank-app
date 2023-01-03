package clients

import (
	"context"
	"fmt"

	"github.com/ivanfomichev/bank-app/internal/database"
	"github.com/jmoiron/sqlx"
)

// Client is a service object
type Client struct {
	Db *sqlx.DB
}

// AddBankClient is a service method to create client
func (c *Client) AddBankClient(ctx context.Context,
	request *database.BankClient,
) (*database.BankClient, error) {
	// logger := applog.FromContext(ctx)

	err := database.AddNewBankClient(ctx, c.Db, request)
	if err != nil {
		fmt.Println("create bank client failed")
		return nil, err
	}
	return &database.BankClient{
		ID: request.ID,
	}, nil
}
