package webapi

import (
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/ivanfomichev/bank-app/internal/database"
)

func (env *RouteHandlers) PostBankClient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := new(database.BankClient)
	err := readValidateInput(ctx, r.Body, req)
	if err != nil {
		log.Printf("bad input")
		BadInputResponse(ctx, w, "create bank client failed")
		return
	}
	err = env.dbclient.GetBankClientByIdentity(ctx, req.IdentityField)
	if err != nil {
		log.Printf("client already exists")
		err = errors.New("client already exists")
		BadInputResponse(ctx, w, err.Error())
	}
	req.ClientID = uuid.New()
	err = env.dbclient.AddBankClient(ctx, req)
	if err != nil {
		log.Printf("create client failed")
		InternalErrorResponse(ctx, w, "create client failed")
		return
	}

	OKResponse(ctx, w, req)
}
