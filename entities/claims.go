package entities

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Role   string `json:"role"`
	Email  string `json:"email"`
	Status string `json:"status"`
	jwt.StandardClaims
}

func (c Claims) IsUserAdmin() bool {
	return strings.EqualFold(c.Role, "admin")
}
func (c Claims) GenerateJwt() (map[string]string, error) {
	secret := os.Getenv("TOKEN_SECRET")
	if secret == "" {
		return nil, errors.New("token secret has not been set")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	signedToken, error := token.SignedString([]byte(secret))
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
	return map[string]string{"access_token": signedToken, "refresh_token": signedRefreshToken}, nil
}
func BuildMapClaims(mapClaims jwt.MapClaims) (*Claims, error) {
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
