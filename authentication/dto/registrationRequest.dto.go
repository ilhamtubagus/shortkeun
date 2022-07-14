package dto

import usr "github.com/ilhamtubagus/urlShortener/user"

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

func (registrationRequestBody *RegistrationRequestBody) ConvertToEntity() *usr.User {
	return &usr.User{
		Name:     registrationRequestBody.Name,
		Email:    registrationRequestBody.Email,
		Password: registrationRequestBody.Password,
		Role:     usr.MEMBER,
		Status:   usr.SUSPENDED,
	}
}
