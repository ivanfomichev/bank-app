// Package webapi - router + handlers
package webapi

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/ivanfomichev/bank-app/internal/clients"
	"github.com/ivanfomichev/bank-app/internal/config"
)

// API contains settings for the web api
type API struct {
	conf           *config.WebAPI
	internalServer *http.Server
}

// RouteHandlers contains env settings for router's handler
type RouteHandlers struct {
	dbclient *clients.Client
}

// APIOption is a type for constructor
type APIOption func(h *RouteHandlers)

// WithDBClient sets database Client for route handlers
func WithDBClient(client *clients.Client) APIOption {
	return func(h *RouteHandlers) {
		h.dbclient = client
	}
}

// NewAPI is a constructor for the *API
func NewAPI(conf *config.WebAPI, options ...APIOption) *API {
	handlerEnv := &RouteHandlers{}
	for i := range options {
		options[i](handlerEnv)
	}
	internalServer := &http.Server{
		Addr:              conf.InternalAPI.Addr,
		Handler:           newInternalRoutes(handlerEnv, conf.InternalAPI),
		ReadHeaderTimeout: 0,
	}

	wapi := &API{
		conf:           conf,
		internalServer: internalServer,
	}
	return wapi
}

// Start launches http-web (REST) API
func (api *API) Start(errCh chan<- error) {

	go func(server *http.Server) {
		log.Printf("starting internal web api")
		err := server.ListenAndServe()
		if err != nil {
			log.Printf("listen to the web api failed")
			errCh <- err
		}
	}(api.internalServer)
}

// Stop terminates http-web (REST) API
func (api *API) Stop(ctx context.Context) error {
	err := api.internalServer.Shutdown(ctx)
	if err != nil {
		log.Printf("failed stop internal server")
		return err
	}
	return nil
}

// nolint:funlen
func newInternalRoutes(env *RouteHandlers, webAPIConf *config.APIConf) http.Handler {
	router := chi.NewRouter()
	router.Use(
		withCORS(webAPIConf.Cors),
	)
	router.Group(func(clientsGroup chi.Router) {
		clientsGroup.Route("/clients", func(clientsRouter chi.Router) {
			clientsRouter.Post("/", env.PostBankClient)
			clientsRouter.Route("/{client_id}", func(clRouter chi.Router) {
				clRouter.Get("/", env.GetBankClientByID)
				clRouter.Post("/", env.PostAccount)
			})
		})
	})
	router.Group(func(accountsGroup chi.Router) {
		accountsGroup.Route("/accounts/{account_id}", func(accountRouter chi.Router) {
			accountRouter.Get("/", env.GetAccountByID)
		})
	})
	router.Group(func(transactionsGroup chi.Router) {
		transactionsGroup.Route("/transactions", func(trRouter chi.Router) {
			trRouter.Get("/", env.GetTransactions)
			// trRouter.Post("/", env.PostTransaction)
		})
	})
	return router
}

func withCORS(conf *config.Cors) func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   conf.AllowedOrigins,
		AllowedMethods:   conf.AllowedMethods,
		AllowedHeaders:   conf.AllowedHeaders,
		Debug:            conf.Debug,
		AllowCredentials: conf.AllowCreds,
	})
}
