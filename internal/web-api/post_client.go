package webapi

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/ivanfomichev/bank-app/internal/database"
)

// BankClient - DTO for client table
type PostBankClient struct {
	ClientID      uuid.UUID `json:"client_id"`
	IdentityField int32     `json:"identity_field" validate:"required"`
	ClientName    string    `json:"client_name" validate:"required"`
}

func (env *RouteHandlers) PostBankClient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := new(PostBankClient)
	err := readValidateInput(ctx, r.Body, req)
	if err != nil {
		log.Printf("bad input")
		BadInputResponse(ctx, w, "create bank client failed")
		return
	}
	req.ClientID = uuid.New()
	err = env.dbclient.AddBankClient(ctx, &database.BankClient{
		ClientID:      req.ClientID,
		IdentityField: req.IdentityField,
		ClientName:    req.ClientName,
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			log.Printf("client already exists")
			err = errors.New("client already exists")
			BadInputResponse(ctx, w, err.Error())
			return
		} else {
			log.Printf("failed to get client from database")
			err = errors.New("failed to get client from database")
			InternalErrorResponse(ctx, w, err.Error())
			return
		}
	}

	OKResponse(ctx, w, req)
}
