package user

import (
	"time"

	"github.com/ilhamtubagus/urlShortener/user/dto"
	"github.com/kamva/mgm/v3"
)

const (
	ADMIN  = "ADMINISTRATOR"
	MEMBER = "MEMBER"
)

const (
	ACTIVE    = "ACTIVE"
	INACTIVE  = "INACTIVE"
	SUSPENDED = "SUSPENDED"
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

func (user User) ConvertToDto() *dto.UserResponseBody {
	return &dto.UserResponseBody{
		ID:     user.ID.String(),
		Email:  user.Email,
		Role:   user.Role,
		Status: user.Status,
		Name:   user.Name,
	}
}

type ActivationCode struct {
	Code     string    `json:"code" bson:"code"`
	IssuedAt time.Time `json:"issued_at" bson:"issued_at"`
	ExpireAt time.Time `json:"expireAt" bson:"expireAt"`
}
