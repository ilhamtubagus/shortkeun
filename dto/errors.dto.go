package dto

import "github.com/ilhamtubagus/urlShortener/lib"

type ValidationErrorResponse struct {
	Message string                 `json:"message"`
	Errors  *[]lib.ValidationError `json:"errors"`
}
