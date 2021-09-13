package dto

type SignInRequestGoogle struct {
	Code string `json:"code" validate:"required"`
}
type SignInRequestDefault struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" valdate:"required,alphanum,min=8,max=25"`
}
type SignInResponse struct {
	Token string `json:"token"`
}

type RegisterRequest struct {
	Name            string `json:"name" validate:"required,max=30"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,alphanum,min=8,max=25"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}
