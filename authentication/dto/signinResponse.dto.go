package dto

// A response when user's successfully signed in
// swagger:response signInResponse
type SignInResponse struct {
	// in: body
	Body SignInResponseBody
}

// swagger:model
type SignInResponseBody struct {
	// The response message
	// Example : signin succeeded
	Message string `json:"message"`
	//	The signin token
	Token TokenResponse `json:"token"`
}

// swagger:model
type TokenResponse struct {
	// access token
	AccessToken string `json:"access_token"`
	// refresh token
	RefreshToken string `json:"refresh_token"`
}
