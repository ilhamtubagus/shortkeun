package dto

// swagger:parameters accountActivation
type AccountActivationRequest struct {
	// in: body
	// required: true
	Body AccountActivationRequestBody
}

// swagger:model
type AccountActivationRequestBody struct {
	// activation code obtained from registration process (sent via email)
	// required: true
	ActivationCode string `json:"activation_code" validate:"required"`
}
