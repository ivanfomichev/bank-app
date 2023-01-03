package webapi

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/ivanfomichev/bank-app/internal/database"
)

// BankClientsResponse is a DTO for BankClients description
type BankClientsResponse struct {
	ID uuid.UUID `json:"id"`
}

func (env *RouteHandlers) PostBankClient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	request := &database.BankClient{
		ID: uuid.New(),
	}
	bankClient, err := env.dbclient.AddBankClient(ctx, request)
	if err != nil {
		log.Printf("create client failed")
		InternalErrorResponse(ctx, w, "create client failed")
		return
	}
	response := &BankClientsResponse{
		ID: bankClient.ID,
	}
	OKResponse(ctx, w, response)
}
