package clients

import (
	"context"
	"log"

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
	err := database.AddNewBankClient(ctx, c.Db, request)
	if err != nil {
		log.Printf("create bank client failed")
		return nil, err
	}
	return &database.BankClient{
		ID: request.ID,
	}, nil
}

// GetBankClient is a service method to get client
func (c *Client) GetBankClient(ctx context.Context,
	clientID string,
) (*database.BankClient, error) {
	bankClient, err := database.GetBankClientByID(ctx, c.Db, clientID)
	if err != nil {
		log.Printf("get bank client failed")
		return nil, err
	}
	return bankClient, nil
}

// AddAccount is a service method to create client
func (c *Client) AddAccount(ctx context.Context,
	request *database.Account,
) (*database.Account, error) {
	err := database.AddNewAccount(ctx, c.Db, request)
	if err != nil {
		log.Printf("create account for bank client failed")
		return nil, err
	}
	return &database.Account{
		ID: request.ID,
	}, nil
}

// GetAccount is a service method to get account
func (c *Client) GetAccountByID(ctx context.Context,
	accountID string,
) (*database.Account, error) {
	account, err := database.GetAccountByID(ctx, c.Db, accountID)
	if err != nil {
		log.Printf("get account failed")
		return nil, err
	}
	return account, nil
}

// GetAccount is a service method to get account
func (c *Client) GetTransactions(ctx context.Context) ([]*database.Transaction, error) {
	transactions, err := database.GetTransactions(ctx, c.Db)
	if err != nil {
		log.Printf("get transactions failed")
		return nil, err
	}
	return transactions, nil
}

// AddTransaction is a service method to create transaction
func (c *Client) AddTransaction(ctx context.Context,
	request *database.Transaction,
) (*database.Transaction, error) {
	if request.TrType == "withdraw" {
		account, err := database.GetAccountByID(ctx, c.Db, request.AccountID.String())
		if err != nil {
			log.Printf("account not found")
			return nil, err
		}
		if account.Balance >= request.Amount:
	}

	
	err := database.AddNewTransaction(ctx, c.Db, request)
	if err != nil {
		log.Printf("create transaction for failed")
		return nil, err
	}
	return &database.Transaction{
		ID: request.ID,
	}, nil
}
