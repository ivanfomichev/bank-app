package webapi

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
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
	response := &BankClientsResponse{
		ID: bankClient.ID,
	}

	OKResponse(ctx, w, response)
}

// func readValidateInput(ctx context.Context, body io.Reader, target interface{}) error {
// 	validate := validator.New()
// 	if err := json.NewDecoder(body).Decode(target); err != nil {
// 		log.Printf("read input failed")
// 		return err
// 	}
// 	if err := validate.Struct(target); err != nil {
// 		log.Printf("validate input failed")
// 		return err
// 	}
// 	return nil
// }
