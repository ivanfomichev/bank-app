package clients

import (
	"context"
	"errors"
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
) error {
	err := database.AddNewBankClient(ctx, c.Db, request)
	if err != nil {
		log.Printf("create bank client failed")
		return err
	}
	return nil
}

// GetBankClient is a service method to get client
func (c *Client) GetBankClientByID(ctx context.Context,
	clientID string,
) (*database.BankClient, error) {
	bankClient, err := database.GetBankClientByID(ctx, c.Db, clientID)
	if err != nil {
		log.Printf("get bank client failed")
		return nil, err
	}
	return bankClient, nil
}

// GetBankClientByIdentity is a service method to get client
func (c *Client) GetBankClientByIdentity(ctx context.Context,
	identityField int32,
) error {
	err := database.GetBankClientByIdentity(ctx, c.Db, identityField)
	if err != nil {
		log.Printf("get bank client failed")
		return err
	}
	return nil
}

// AddAccount is a service method to create account
func (c *Client) AddAccount(ctx context.Context,
	request *database.Account,
) error {
	err := database.AddNewAccount(ctx, c.Db, request)
	if err != nil {
		log.Printf("create account for bank client failed")
		return err
	}
	return  nil
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

// GetTransactions is a service method to get transactions
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
) error {
	switch recType := request.TrType; recType {
	case "withdraw":
		{
			account, err := database.GetAccountByID(ctx, c.Db, request.AccountID.String())
			if err != nil {
				log.Printf("account not found")
				return err
			}
			if account.Balance >= request.Amount {
				tx, err := c.Db.BeginTxx(ctx, nil)
				if err != nil {
					log.Printf("failed to start db_transaction")
					return err
				}
				err = database.AddNewTransaction(ctx, tx, request)
				if err != nil {
					tx.Rollback()
					return err
				}
				newBalance := account.Balance - request.Amount
				err = database.UpdateAccountByID(ctx, tx, request.AccountID, newBalance)
				if err != nil {
					tx.Rollback()
					return err
				}
				err = tx.Commit()
				if err != nil {
					log.Printf("failed to commit db_transaction")
					return err
				}
				return nil
			} else {
				log.Printf("not enough money")
				err = errors.New("not enough money")
				return err
			}
		}
	case "deposit":
		{
			account, err := database.GetAccountByID(ctx, c.Db, request.AccountID.String())
			if err != nil {
				log.Printf("account not found")
				return err
			}
			tx, err := c.Db.BeginTxx(ctx, nil)
			if err != nil {
				log.Printf("failed to start db_transaction")
				return err
			}
			err = database.AddNewTransaction(ctx, tx, request)
			if err != nil {
				tx.Rollback()
				return err
			}
			newBalance := account.Balance + request.Amount
			err = database.UpdateAccountByID(ctx, tx, request.AccountID, newBalance)
			if err != nil {
				tx.Rollback()
				return err
			}
			err = tx.Commit()
			if err != nil {
				log.Printf("failed to commit db_transaction")
				return err
			}
			return nil
		}
	case "transfer":
		{
			reqId := request.AccountID.String()
			accountFrom, err := database.GetAccountByID(ctx, c.Db, reqId)
			if err != nil {
				log.Printf("from account not found")
				return err
			}
			accountTo, err := database.GetAccountByID(ctx, c.Db, request.AccountToID.String())
			if err != nil {
				log.Printf("to account not found")
				return err
			}
			if accountFrom.Currency != accountTo.Currency {
				log.Printf("transaction for different currencies not allowed")
				err = errors.New("transaction for different currencies not allowed")
				return err
			}
			if request.AccountID == request.AccountToID {
				return nil
			}
			if accountFrom.Balance >= request.Amount {
				tx, err := c.Db.BeginTxx(ctx, nil)
				if err != nil {
					log.Printf("failed to start db_transaction")
					return err
				}
				err = database.AddNewTransaction(ctx, tx, request)
				if err != nil {
					tx.Rollback()
					return err
				}
				newToBalance := accountTo.Balance + request.Amount
				err = database.UpdateAccountByID(ctx, tx, request.AccountToID, newToBalance)
				if err != nil {
					tx.Rollback()
					return err
				}
				newFromBalance := accountFrom.Balance - request.Amount
				err = database.UpdateAccountByID(ctx, tx, request.AccountID, newFromBalance)
				if err != nil {
					tx.Rollback()
					return err
				}
				err = tx.Commit()
				if err != nil {
					log.Printf("failed to commit db_transaction")
					return err
				}
				return nil
			} else {
				log.Printf("not enough money")
				err = errors.New("not enough money")
				return err
			}
		}
	}
	return nil
}
