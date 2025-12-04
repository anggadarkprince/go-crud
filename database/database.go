package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/anggadarkprince/crud-employee-go/configs"
	_ "github.com/go-sql-driver/mysql"
)

func InitDatabase() *sql.DB {
	db, err := sql.Open("mysql", configs.Get().Database.DSN())
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(10)                 // Max number of open connections
	db.SetMaxIdleConns(5)                  // Max number of idle connections
	db.SetConnMaxLifetime(5 * time.Minute) // Max lifetime for a connection

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

type Transaction interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}
