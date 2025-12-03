package exceptions

import "fmt"

type AppError struct {
	Code int
	Message string
	Err error
}

func (e *AppError) Error() string {
	return fmt.Sprintf("Code %d: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
    return e.Err
}