package webapi

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/ivanfomichev/bank-app/internal/database"
)

// AccountsResponse is a DTO for Accounts description
type AccountResponse struct {
	ID uuid.UUID `json:"id"`
}

func (env *RouteHandlers) PostAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	clID := chi.URLParam(r, "client_id")
	request := &database.Account{
		ID:           uuid.New(),
		BankClientID: clID,
	}
	account, err := env.dbclient.AddAccount(ctx, request)
	if err != nil {
		log.Printf("create client failed")
		InternalErrorResponse(ctx, w, "create client failed")
		return
	}
	response := &AccountResponse{
		ID: account.ID,
	}
	OKResponse(ctx, w, response)
}
