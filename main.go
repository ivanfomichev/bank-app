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

// Hello Ivan. Thank you for applying for the position. However, the test met the minimum requirements to pass
// Here are some comments about the test.
// 1. SQL schema
// - ‘accounts’ table does refer to bank_clients
// - balance can be negative

// 2. Create client (post request).
// - ‘Client already exists’ error is not represented

// 3. Create Account (post request).
// - Account with invalid client ID can be created
// - ‘Invalid currency’ error is not presented

// 4. Add transaction
// -Withdraw. Big chance to get inconsistent balance because account is not locked
// - Transfer. Big chance to get inconsistent balance because account is not locked