package utilities

import (
	"database/sql"
	"time"
)

func StringToDate(value string) (sql.NullTime, error) {
	var formatterValue sql.NullTime;
	if value != "" {
		t, err := time.Parse("2006-01-02", value)
		if err != nil {
			return formatterValue, err
		}

		formatterValue = sql.NullTime{
			Time:  t,
			Valid: true,
		}
	}
	return formatterValue, nil
}
