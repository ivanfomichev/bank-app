// Package main package
package main

import (
	"context"
	"log"

	"github.com/ivanfomichev/bank-app/internal/app"
	conf "github.com/ivanfomichev/bank-app/internal/config"
)

func main() {
	ctx := context.Background()
	config, err := conf.New()
	if err != nil {
		log.Fatal(err)
	}

	stopMe, errCh := app.Start(ctx, config)

	err, ok := <-errCh
	if ok {
		if err != nil {
			log.Printf("application failed")
		}
	} else {
		log.Printf("error chan closed")
	}
	stopMe()
}
