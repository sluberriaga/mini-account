package account

import "fmt"

type validationError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewValidationError(code, field string) validationError {
	return validationError{
		code,
		fmt.Sprintf("field %s has errors!", field),
	}
}
