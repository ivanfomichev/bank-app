package webapi

import (
	"net/http"
)

// type PostTransactionRequest struct {
// 	ID          uuid.UUID `json:"id" validate:"required"`
// 	AccountID   uuid.UUID `json:"account_id" validate:"required"`
// 	AccountToID uuid.UUID `json:"account_to_id"`
// 	TrType      string    `json:"tr_type" validate:"required"`
// 	TrStatus    string    `json:"tr_status"`
// }

func (env *RouteHandlers) PostTransaction(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()

	// trType := chi.URLParam(r, "tr_type")

	// req := new(database.Transaction)
	// err := readValidateInput(ctx, r.Body, req)
	// if err != nil {
	// 	log.Printf("bad input")
	// 	BadInputResponse(ctx, w, "create account failed")
	// }
	// req.ID = uuid.New()

	// account, err := env.dbclient.AddTransaction(ctx, req)
	// if err != nil {
	// 	log.Printf("create account failed")
	// 	InternalErrorResponse(ctx, w, "create account failed")
	// 	return
	// }
	// OKResponse(ctx, w, PostResponse{
	// 	ID: account.ID,
	// })
}
