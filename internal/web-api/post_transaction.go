package webapi

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/ivanfomichev/bank-app/internal/database"
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
	req.ID = uuid.New()
	account, err := env.dbclient.AddTransaction(ctx, req)
	if err != nil {
		log.Printf("create transaction failed")
		InternalErrorResponse(ctx, w, err.Error())
		return
	}
	OKResponse(ctx, w, PostResponse{
		ID: account.ID,
	})
}
