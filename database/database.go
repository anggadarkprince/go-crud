package database

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func InitDatabase() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/sandbox")
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(10) // Max number of open connections
	db.SetMaxIdleConns(5) // Max number of idle connections
	db.SetConnMaxLifetime(5 * time.Minute)  // Max lifetime for a connection

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}