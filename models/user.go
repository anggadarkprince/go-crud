package models

import "database/sql"

type User struct {
	Id int
	Name string
	Username string
	Email string
	Password string
	UserType string
	Status string
	Avatar sql.NullString
}