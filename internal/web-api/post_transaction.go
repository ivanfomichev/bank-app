package webapi

import (
	"errors"
	"log"
	"net/http"
	"strings"

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
	_, err = env.dbclient.GetAccountByID(ctx, req.AccountID.String())
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			log.Printf("no such account")
			err = errors.New("no accounts with specified account_id")
			BadInputResponse(ctx, w, err.Error())
			return
		} else {
			log.Printf("failed to get account from database")
			BadInputResponse(ctx, w, "create transaction failed")
			return
		}
	}

	if req.TrType == Transfer {
		_, err := env.dbclient.GetAccountByID(ctx, req.AccountToID.String())
		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				log.Printf("no such account")
				err = errors.New("no accounts with specified account_id")
				BadInputResponse(ctx, w, err.Error())
				return
			} else {
				log.Printf("failed to get account from database")
				BadInputResponse(ctx, w, "create transaction failed")
				return
			}
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
