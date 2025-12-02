package models

import "database/sql"

type Employee struct {
	Id int
	Name string
	Email sql.NullString
	TaxNumber sql.NullString
	Gender sql.NullString
	HiredDate sql.NullTime
	Address sql.NullString
	Status sql.NullString
	TotalAllowance int
}