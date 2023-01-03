// Package app - main application package
package app

import (
	"context"
	"fmt"

	"github.com/ivanfomichev/bank-app/internal/clients"
	conf "github.com/ivanfomichev/bank-app/internal/config"
	"github.com/ivanfomichev/bank-app/internal/database"

	webapi "github.com/ivanfomichev/bank-app/internal/web-api"
)

// appServices - application main engine
type appServices struct {
	database *clients.Client
}

// Start - start operation
func Start(ctx context.Context, conf *conf.Config) (StopFunc, <-chan error) {
	errCh := make(chan error, 1)

	appsrv, err := initAppServices(ctx, conf)
	if err != nil {
		fmt.Println("failed init application services")
		errCh <- err
		return func() {}, errCh
	}
	options := []webapi.APIOption{
		webapi.WithDBClient(appsrv.database),
	}
	webAPI := webapi.NewAPI(conf.WebAPI, options...)

	go func() {
		webAPI.Start(errCh)
	}()

	stop := StopFunc(
		func() {
			if stopErr := webAPI.Stop(ctx); stopErr != nil {
				fmt.Println("stop web api failed")
			}
		})
	return stop, errCh
}

// StopFunc is a application terminating func
// should be used by app.Start caller
type StopFunc func()

// Stop - graceful finishing function
func (a *appServices) Stop() {
	fmt.Printf("Application has been gracefully stopped")
}

func initAppServices(ctx context.Context, conf *conf.Config) (*appServices, error) {
	dbClient, err := database.OpenDatabase(ctx, conf.Database, database.Reg())
	if err != nil {
		fmt.Println("create database client failed")
		return nil, err
	}
	return &appServices{
		database: &clients.Client{
			Db: dbClient,
		},
	}, nil
}
