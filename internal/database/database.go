// Package database - application database connector
package database

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"log"

	"github.com/google/uuid"
	pgDriver "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/luna-duclos/instrumentedsql"
	"github.com/luna-duclos/instrumentedsql/opentracing"

	"github.com/ivanfomichev/bank-app/internal/config"
)

const (
	wrappedDriverName = "cloudsqlpostgres"
)

// SQLExecutor this interface work with *sqlx.DB and *sqlx.Tx structs and recommended for work with database
type SQLExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

// SQLRebinder interface for enabling rebinding dollars and question marks in SQL queries
type SQLRebinder interface {
	Rebind(query string) string
}

// SQLExecRebinder union interface of SQLExecutor and SQLRebinder
type SQLExecRebinder interface {
	SQLExecutor
	SQLRebinder
}

// Reg register new wrapedDriver for trace
func Reg() string {
	drv := wrappedDriverName
	if !isDriverRegistered(drv) {
		sql.Register(drv,
			instrumentedsql.WrapDriver(&pgDriver.Driver{},
				instrumentedsql.WithTracer(opentracing.NewTracer(true)),
			),
		)
	}
	return wrappedDriverName
}

// OpenDatabase open postgresql database from DbConnectString
func OpenDatabase(ctx context.Context, conf *config.DatabaseConfig, driverName string) (*sqlx.DB, error) {
	conn, err := sqlx.ConnectContext(ctx, driverName, conf.ConnString)
	if err != nil {
		log.Printf("Connection to database failed")
		return nil, err
	}
	return conn, nil
}

// CloseDatabase close db database
func CloseDatabase(ctx context.Context, db io.Closer) error {
	err := db.Close()
	if err != nil {
		log.Printf("Close connection to database failed")
	}
	return err
}

func isDriverRegistered(drv string) bool {
	driversList := sql.Drivers()
	for _, driverName := range driversList {
		if driverName == drv {
			return true
		}
	}
	return false
}

func execInsertObjectQuery(ctx context.Context, dbc SQLExecutor, query string, cs interface{}) error {
	if cs == nil {
		return errors.New("no check provided")
	}
	_, err := dbc.NamedExecContext(ctx, query, cs)
	if err != nil {
		log.Printf("failed add new object")
	}
	return err
}

func updateTableColWithProvidedKey(
	ctx context.Context, dbc sqlx.ExecerContext, query string, colValue interface{}, keyID uuid.UUID,
) error {
	r, err := dbc.ExecContext(ctx, query, colValue, keyID)
	if err != nil {
		log.Printf("failed modify row")
		return err
	}
	affected, err := r.RowsAffected()
	if err != nil {
		log.Printf("cannot get results of update")
		return nil
	}
	if affected == 0 {
		log.Printf("unknown id")
		return sql.ErrNoRows
	}
	return err
}

func scanTableRowToObject(
	ctx context.Context, dbc sqlx.QueryerContext, query string, obj interface{}, keyID uuid.UUID,
) error {
	if obj == nil {
		return errors.New("no object provided")
	}
	err := dbc.QueryRowxContext(ctx, query, keyID).StructScan(obj)
	if err != nil {
		log.Printf("failed find object")
	}
	return err
}

func deleteAllRowsWithProvidedKey(
	ctx context.Context, dbc sqlx.ExecerContext, query string, keyID uuid.UUID,
) (int64, error) {
	r, err := dbc.ExecContext(ctx, query, keyID)
	if err != nil {
		log.Printf("failed remove rows")
		return 0, err
	}
	affected, err := r.RowsAffected()
	if err != nil {
		log.Printf("cannot get results of update")
		return 0, nil
	}
	if affected == 0 {
		log.Printf("unknown id")
		return affected, sql.ErrNoRows
	}
	return affected, err
}
