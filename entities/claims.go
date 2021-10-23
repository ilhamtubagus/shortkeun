package entities

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// swagger:model
type Token struct {
	// access token
	AccessToken string `json:"access_token"`
	// refresh token
	RefreshToken string `json:"refresh_token"`
}
type Claims struct {
	UserId string `json:"userId"`
	Role   string `json:"role"`
	Email  string `json:"email"`
	Status string `json:"status"`
	jwt.StandardClaims
}

func (c Claims) IsUserAdmin() bool {
	return strings.EqualFold(c.Role, "admin")
}
func (c Claims) GenerateJwt() (*Token, error) {
	secret := os.Getenv("TOKEN_SECRET")
	if secret == "" {
		return nil, errors.New("token secret has not been set")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	signedAccessToken, error := token.SignedString([]byte(secret))
	if error != nil {
		return nil, error
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = c.Subject
	//refresh token expires within 24 hours
	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	signedRefreshToken, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}
	return &Token{AccessToken: signedAccessToken, RefreshToken: signedRefreshToken}, nil
}
func BuildMapClaims(mapClaims jwt.Claims) (*Claims, error) {
	bytes, errs := json.Marshal(mapClaims)
	if errs != nil {
		return nil, errs
	}
	var claims Claims
	error := json.Unmarshal(bytes, &claims)
	if error != nil {
		return nil, error
	}
	return &claims, nil
}
