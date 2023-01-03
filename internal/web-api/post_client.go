package webapi

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/ivanfomichev/bank-app/internal/database"
)

// BankClientData is input DTO for POST /clients
type BankClientData struct {
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	BirthDate string `json:"birth_date"`
	Passport  string `json:"passport"`
	Job       string `json:"job"`
}

// BankClientsResponse is a DTO for BankClients description
type BankClientsResponse struct {
	ID uuid.UUID `json:"id"`
}

// PostClients godoc
// @Accept json
// @Produce json
// @Success 200 {object} ResponseBody{data=BankClientsDescResponse} "desc"
// @Failure 400 {object} ResponseBody{} "desc"
// @Router /clients [post]
//
//nolint:funlen
func (env *RouteHandlers) PostBankClient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	bankClientData := new(BankClientData)
	if err := readValidateInput(ctx, r.Body, bankClientData); err != nil {
		log.Printf("failed read or validate input body")
		BadInputResponse(ctx, w, "invalid input body")
		return
	}
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
