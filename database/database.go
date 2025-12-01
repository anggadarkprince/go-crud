package database

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func InitDatabase() *sql.DB {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_DATABASE")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")

	db, err := sql.Open("mysql", username + ":" + password + "@tcp(" + host + ":" + port + ")/" + database + "?parseTime=true")
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