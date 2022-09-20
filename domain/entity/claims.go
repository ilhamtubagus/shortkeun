package entity

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
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
	refreshTokenExpirationInString := os.Getenv("REFRESH_TOKEN_EXP")
	if refreshTokenExpirationInString == "" {
		return nil, errors.New("refresh token expiration has not been set")
	}
	hour, err := strconv.Atoi(refreshTokenExpirationInString)
	if err != nil {
		return nil, err
	}
	rtClaims["exp"] = time.Now().Add(time.Hour * time.Duration(hour)).Unix()

	signedRefreshToken, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}
	return &Token{AccessToken: signedAccessToken, RefreshToken: signedRefreshToken}, nil
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
