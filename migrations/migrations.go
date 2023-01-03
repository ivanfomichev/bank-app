// Package migrations - database migrations utilities
package migrations

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	migrate "github.com/rubenv/sql-migrate"

	"github.com/ivanfomichev/bank-app/internal/config"
	"github.com/ivanfomichev/bank-app/internal/database"
)

// Options - database migration options
type Options struct {
	conf      *config.Config
	count     int
	direction migrate.MigrationDirection
}

const defaultSchemaName = "public"

// New - migrations options constructor
func New(conf *config.Config, direction migrate.MigrationDirection, count int) *Options {
	return &Options{
		conf:      conf,
		count:     count,
		direction: direction,
	}
}

// Run - start new migrations
func (o *Options) Run(ctx context.Context) error {

	log.Printf("Applying migrations")

	dbConn, err := database.OpenDatabase(ctx, o.conf.Database, database.Reg())
	if err != nil {
		log.Printf("db connect failed")
	}
	defer func() {
		clErr := dbConn.Close()
		if clErr != nil {
			log.Printf("closing db connection is failed")
		}
	}()

	log.Printf("Migrations destination is: %v", o.conf.Database.ConnString)

	dbConf, err := pgx.ParseConfig(o.conf.Database.ConnString)
	if err != nil {
		log.Printf("failed parse database connection string")
	}
	schemaName, ok := dbConf.RuntimeParams["search_path"]
	if !ok {
		schemaName = defaultSchemaName
	}
	migrate.SetSchema(schemaName)
	count, err := migrate.ExecMax(dbConn.DB, "postgres",
		&migrate.FileMigrationSource{Dir: "./migrations/sql"},
		o.direction, o.count,
	)
	if err != nil {
		log.Println(err.Error())
	}
	log.Printf("Applied %v migrations", count)
	return err
}
