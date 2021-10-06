package dto

import "github.com/ilhamtubagus/urlShortener/entities"

// Request schema for default sign in
// swagger:parameters signIn
type SignInRequestDefault struct {
	//
	// in: body
	// required: true
	Body SignInRequestDefaultBody
}
type SignInRequestDefaultBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,alphanum,min=8,max=25"`
}

// Request schema for sign in with google account
// swagger:parameters googleSignIn
type SignInRequestGoogle struct {
	//
	// in: body
	// required: true
	Body GoogleSignInRequestBody
}

// swagger:model
type GoogleSignInRequestBody struct {
	// contain JWT ID Token obtained from google
	Credential string `json:"credential" validate:"required"`
}

// swagger:parameters register
type RegistrationRequest struct {
	// in: body
	// required: true
	Body RegistrationRequestBody
}

// swagger:model
type RegistrationRequestBody struct {
	// users fullname
	// required: true
	// max length: 30
	Name string `json:"name" validate:"required,max=30"`
	// email
	// required: true
	// swagger:strfmt email
	Email string `json:"email" validate:"required,email"`
	// password
	// required: true
	// min length: 8
	// max length: 25
	Password string `json:"password" validate:"required,alphanum,min=8,max=25"`
	// must be equal with password
	// required: true
	// min length: 8
	// max length: 25
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

// swagger:parameters getActivationCode
type ActivationCodeRequest struct {
	// in: body
	// required: true
	Body ActivationCodeRequestBody
}

//	swagger: model
type ActivationCodeRequestBody struct {
	// email
	// required: true
	// swagger:strfmt email
	Email string `json:"email" validate:"required,email"`
}

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
	Token entities.Token `json:"token"`
}
