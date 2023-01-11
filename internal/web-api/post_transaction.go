package webapi

import (
	"errors"
	"log"
	"net/http"
	"reflect"

	"github.com/google/uuid"
	"github.com/ivanfomichev/bank-app/internal/database"
)

const (
	Withdrow = "withdrow"
	Deposit  = "deposit"
	Transfer = "transfer"
)

func (env *RouteHandlers) PostTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := new(database.Transaction)
	err := readValidateInput(ctx, r.Body, req)
	if err != nil {
		log.Printf("bad input")
		BadInputResponse(ctx, w, "create transaction failed")
		return
	}
	account, err := env.dbclient.GetAccountByID(ctx, req.AccountID.String())
	if err != nil {
		log.Printf("failed to get account from database")
		BadInputResponse(ctx, w, "create transaction failed")
	}
	if reflect.DeepEqual(account, &database.Account{}) {
		log.Printf("no such account")
		err = errors.New("no account with specified account_id")
		BadInputResponse(ctx, w, err.Error())
	}

	if req.TrType == Transfer {
		account, err := env.dbclient.GetAccountByID(ctx, req.AccountToID.String())
		if err != nil {
			log.Printf("failed to get account from database")
			BadInputResponse(ctx, w, "create transaction failed")
		}
		if reflect.DeepEqual(account, &database.Account{}) {
			log.Printf("no such account")
			err = errors.New("no account with specified account_id")
			BadInputResponse(ctx, w, err.Error())
		}
	}
	req.TransactionID = uuid.New()
	err = env.dbclient.AddTransaction(ctx, req)
	if err != nil {
		log.Printf("create transaction failed")
		InternalErrorResponse(ctx, w, err.Error())
		return
	}
	OKResponse(ctx, w, req)
}
