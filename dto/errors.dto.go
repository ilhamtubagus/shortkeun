package dto

import "github.com/ilhamtubagus/urlShortener/lib"

// A ValidationError is an error that is used when the required input fails validation.
// swagger:response validationError
type ValidationErrorResponse struct {
	// The message
	// in:body
	Message string `json:"message"`
	// The list of error
	// in:body
	Errors *[]lib.ValidationError `json:"errors"`
}
