package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/ilhamtubagus/urlShortener/domain/entity"
)

func getGooglePublicKey(keyID string) (string, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v1/certs")
	if err != nil {
		return "", err
	}
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	myResp := map[string]string{}
	err = json.Unmarshal(dat, &myResp)
	if err != nil {
		return "", err
	}
	key, ok := myResp[keyID]
	if !ok {
		return "", errors.New("key not found")
	}
	return key, nil
}

type Oauth2GoogleService interface {
	VerifyToken(idToken string) (*entity.GoogleClaims, error)
}

type oauth2GoogleService struct {
	googleClientId string
}

func (o oauth2GoogleService) VerifyToken(idToken string) (*entity.GoogleClaims, error) {
	googleClaims := entity.GoogleClaims{}
	token, err := jwt.ParseWithClaims(idToken, &googleClaims, func(t *jwt.Token) (interface{}, error) {
		pem, err := getGooglePublicKey(fmt.Sprintf("%s", t.Header["kid"]))
		if err != nil {
			return nil, err
		}
		key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
		if err != nil {
			return nil, err
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*entity.GoogleClaims)
	if !ok {
		return nil, errors.New("invalid google jwt")
	}
	if claims.Issuer != "accounts.google.com" && claims.Issuer != "https://accounts.google.com" {
		return nil, errors.New("iss is invalid")
	}

	if claims.Audience != o.googleClientId {
		return nil, errors.New("aud is invalid")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return nil, errors.New("jwt is expired")
	}
	return claims, nil
}

func NewOauth2GoogleService(googleClientId string) Oauth2GoogleService {
	return oauth2GoogleService{googleClientId: googleClientId}
}
