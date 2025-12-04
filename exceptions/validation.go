package exceptions

import "fmt"

type ValidationError struct {
	Message string
	Errors  map[string]string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Validation error: %s", e.Message)
}
