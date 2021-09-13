package dto

type SignInRequestGoogle struct {
	Code string `json:"code" validate:"required"`
}
type SignInRequestDefault struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type SignInResponse struct {
	Token string `json:"token"`
}
