package dto

import "github.com/ilhamtubagus/urlShortener/user"

// swagger:parameters getActivationCode
type NewActivationCodeRequest struct {
	// in: body
	// required: true
	Body NewActivationCodeRequestBody
}

//	swagger: model
type NewActivationCodeRequestBody struct {
	// email
	// required: true
	// swagger:strfmt email
	Email string `json:"email" validate:"required,email"`
}

func (newActivationCodeRequestBody NewActivationCodeRequestBody) ConvertToEntity() *user.User {
	return &user.User{
		Email: newActivationCodeRequestBody.Email,
	}
}
