package webapi

import (
	"errors"
	"log"
	"net/http"
	"reflect"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/ivanfomichev/bank-app/internal/database"
)

var validCurr = map[string]bool{
	"MXN": true,
	"COP": true,
	"USD": true,
}

func (env *RouteHandlers) PostAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	clID := chi.URLParam(r, "client_id")
	// validate bank client
	client, err := env.dbclient.GetBankClientByID(ctx, clID)
	if err != nil {
		log.Printf("failed to get client from database")
		BadInputResponse(ctx, w, "create account failed")
	}
	if reflect.DeepEqual(client, &database.BankClient{}) {
		log.Printf("no such client")
		err = errors.New("no clients with specified client_id")
		BadInputResponse(ctx, w, err.Error())
	}

	req := new(database.Account)
	err = readValidateInput(ctx, r.Body, req)
	if err != nil {
		log.Printf("bad input")
		BadInputResponse(ctx, w, "create account failed")
	}
	// validate currency type
	if ok := validCurr[req.Currency]; !ok {
		log.Printf("currency not valid")
		err = errors.New("currency not valid")
		BadInputResponse(ctx, w, err.Error())
	}

	req.AccountID = uuid.New()
	req.ClientID, err = uuid.Parse(clID)
	if err != nil {
		log.Printf("bad input")
		BadInputResponse(ctx, w, "create account failed")
	}
	err = env.dbclient.AddAccount(ctx, req)
	if err != nil {
		log.Printf("create account failed")
		InternalErrorResponse(ctx, w, "create account failed")
		return
	}

	OKResponse(ctx, w, req)
}
