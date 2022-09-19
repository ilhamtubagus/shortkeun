package dto

// swagger:parameters register
type RefreshTokenRequest struct {
	// in: body
	// required: true
	Body RefreshTokenRequestBody
}
type RefreshTokenRequestBody struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
