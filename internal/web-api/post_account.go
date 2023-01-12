package webapi

import (
	"errors"
	"log"
	"net/http"
	"strings"

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
	req := new(database.Account)
	err := readValidateInput(ctx, r.Body, req)
	if err != nil {
		log.Printf("bad input")
		BadInputResponse(ctx, w, "create account failed")
		return
	}
	// validate bank client
	_, err = env.dbclient.GetBankClientByID(ctx, clID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			log.Printf("no such client")
			err = errors.New("no clients with specified client_id")
			BadInputResponse(ctx, w, err.Error())
			return
		} else {
			log.Printf("failed to get client from database")
			InternalErrorResponse(ctx, w, "create account failed")
			return
		}
	}
	// validate currency type
	if ok := validCurr[req.Currency]; !ok {
		log.Printf("currency not valid")
		err = errors.New("currency not valid")
		BadInputResponse(ctx, w, err.Error())
		return
	}

	req.AccountID = uuid.New()
	uid, err := uuid.Parse(clID)
	if err != nil {
		log.Printf("can not parse uuid from string")
		InternalErrorResponse(ctx, w, "create account failed")
		return
	}
	req.ClientID = uid
	err = env.dbclient.AddAccount(ctx, req)
	if err != nil {
		log.Printf("create account failed")
		InternalErrorResponse(ctx, w, "create account failed")
		return
	}

	OKResponse(ctx, w, req)
}
