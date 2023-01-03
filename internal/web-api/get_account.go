package webapi

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func (env *RouteHandlers) GetAccountByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accID := chi.URLParam(r, "account_id")
	account, err := env.dbclient.GetAccountByID(ctx, accID)
	if err != nil {
		log.Printf("get client failed")
		InternalErrorResponse(ctx, w, "get client failed")
		return
	}

	OKResponse(ctx, w, account)
}
