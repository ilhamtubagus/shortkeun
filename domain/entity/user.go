package entity

import (
	"github.com/ilhamtubagus/urlShortener/interface/dto"
	"github.com/kamva/mgm/v3"
)

type User struct {
	mgm.DefaultModel `bson:",inline"`
	Email            string `json:"email"`
	Name             string `json:"name"`
	Password         string `json:"-" bson:",omitempty"`
	Status           string `json:",omitempty"`
	Role             string `json:"role"`
	// subject from google
	Sub            string          `json:"sub" bson:"sub,omitempty"`
	ActivationCode *ActivationCode `json:"activation_code,omitempty" bson:"activation_code,omitempty"`
}

func (user User) ConvertToResponseDto() *dto.UserResponseBody {
	return &dto.UserResponseBody{
		ID:     user.ID.Hex(),
		Email:  user.Email,
		Role:   user.Role,
		Status: user.Status,
		Name:   user.Name,
	}
}
