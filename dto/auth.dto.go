package dto

type SignInRequestDefault struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,alphanum,min=8,max=25"`
}
type SignInRequestGoogle struct {
	Credential string `json:"credential" validate:"required"`
}
type SignInResponse struct {
	Message string            `json:"message"`
	Token   map[string]string `json:"token"`
}
type RegisterRequest struct {
	Name            string `json:"name" validate:"required,max=30"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,alphanum,min=8,max=25"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type RequestCodeActivation struct {
	Email string `json:"email" validate:"required,email"`
}

type AccountActivationRequest struct {
	ActivationCode string `json:"activation_code" validate:"required"`
}
