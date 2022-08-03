package dto

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
