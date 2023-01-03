package webapi

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

func (env *RouteHandlers) GetBankClientByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	clID := chi.URLParam(r, "client_id")
	bankClient, err := env.dbclient.GetBankClient(ctx, clID)
	if err != nil {
		log.Printf("get client failed")
		InternalErrorResponse(ctx, w, "get client failed")
		return
	}
	response := &PostResponse{
		ID: bankClient.ID,
	}

	OKResponse(ctx, w, response)
}

func readValidateInput(ctx context.Context, body io.Reader, target interface{}) error {
	validate := validator.New()
	if err := json.NewDecoder(body).Decode(target); err != nil {
		log.Printf("read input failed")
		return err
	}
	if err := validate.Struct(target); err != nil {
		log.Printf("validate input failed")
		return err
	}
	return nil
}
