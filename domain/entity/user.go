package entity

import (
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
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

func (user User) CreateClaims() (*Claims, error) {
	//create our own jwt and send back to client
	hour, err := strconv.Atoi(os.Getenv("TOKEN_EXP"))
	if err != nil {
		return nil, err
	}
	claims := Claims{
		Role:   user.Role,
		Email:  user.Email,
		Status: user.Status,
		StandardClaims: jwt.StandardClaims{
			//token expires within x hours
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(hour)).Unix(),
			Subject:   user.ID.String(),
		}}
	return &claims, nil
}
