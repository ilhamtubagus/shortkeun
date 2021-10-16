package dto

import (
	"github.com/ilhamtubagus/urlShortener/lib"
)

// A ValidationError is an error that is used when the required input fails validation.
// swagger:response validationError
type ValidationErrorResponse struct {
	// in:body
	// required: true
	Body ValidationErrorResponseBody
}

type ValidationErrorResponseBody struct {
	// The message
	Message string `json:"message"`
	// Field errors with its messages
	Errors *[]lib.ValidationError `json:"errors"`
}

func NewValidationError(msg string, err *[]lib.ValidationError) *ValidationErrorResponseBody {
	return &ValidationErrorResponseBody{msg, err}
}

type ApiError struct {
	StatusCode int
	Err        error
}

func NewApiError(code int, err error) *ApiError {
	return &ApiError{
		StatusCode: code,
		Err:        err,
	}
}
