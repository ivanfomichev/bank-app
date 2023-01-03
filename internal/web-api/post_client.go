package webapi

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/ivanfomichev/bank-app/internal/database"
)

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
	response := &PostResponse{
		ID: bankClient.ID,
	}
	OKResponse(ctx, w, response)
}
