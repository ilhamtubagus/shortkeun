package entity

import "github.com/dgrijalva/jwt-go"

type GoogleClaims struct {
	Email         string `json:"email"`
	Name          string `json:"name"`
	EmailVerified bool   `json:"email_verified"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
	Sub           string `json:"sub"`
	jwt.StandardClaims
}
