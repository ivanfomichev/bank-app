package webapi

import (
	"log"
	"net/http"
)

func (env *RouteHandlers) GetTransactions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	transactions, err := env.dbclient.GetTransactions(ctx)
	if err != nil {
		log.Printf("get transactions failed")
		InternalErrorResponse(ctx, w, "get transactions failed")
		return
	}

	OKResponse(ctx, w, transactions)
}
