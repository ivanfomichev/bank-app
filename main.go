// Package main package
package main

import (
	"context"
	"log"

	"github.com/ivanfomichev/bank-app/internal/app"
	conf "github.com/ivanfomichev/bank-app/internal/config"
	"github.com/ivanfomichev/bank-app/migrations"
	migrate "github.com/rubenv/sql-migrate"
)

func main() {
	ctx := context.Background()
	config, err := conf.New()
	if err != nil {
		log.Fatal(err)
	}
	
	err = makeMigrations(ctx, config)
	if err != nil {
		log.Fatal("migrations failed")
	}

	stopMe, errCh := app.Start(ctx, config)

	err, ok := <-errCh
	if ok {
		if err != nil {
			log.Fatal("application failed")
		}
	} else {
		log.Printf("error chan closed")
	}
	stopMe()
}

func makeMigrations(ctx context.Context, conf *conf.Config) error {
	return migrations.New(conf, migrate.Up, 0).Run(ctx)
}
