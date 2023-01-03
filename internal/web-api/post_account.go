package webapi

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/ivanfomichev/bank-app/internal/database"
)

// PostResponse is a DTO for BankClients description
type PostResponse struct {
	ID uuid.UUID `json:"id"`
}

func (env *RouteHandlers) PostAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	clID := chi.URLParam(r, "client_id")
	req := new(database.Account)
	err := readValidateInput(ctx, r.Body, req)
	if err != nil {
		log.Printf("bad input")
		InternalErrorResponse(ctx, w, "create account failed")
	}
	req.ID = uuid.New()
	req.BankClientID = clID
	account, err := env.dbclient.AddAccount(ctx, req)
	if err != nil {
		log.Printf("create account failed")
		InternalErrorResponse(ctx, w, "create account failed")
		return
	}
	OKResponse(ctx, w, PostResponse{
		ID: account.ID,
	})
}
