package dto

import "github.com/ilhamtubagus/urlShortener/user"

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

func (signInRequestDefaultBody SignInRequestDefaultBody) ConvertToEntity() *user.User {
	return &user.User{
		Email:    signInRequestDefaultBody.Email,
		Password: signInRequestDefaultBody.Password,
	}
}
