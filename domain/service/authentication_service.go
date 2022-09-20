package service

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/ilhamtubagus/urlShortener/domain/constant"
	"github.com/ilhamtubagus/urlShortener/domain/entity"
	"github.com/ilhamtubagus/urlShortener/interface/dto"
	"github.com/ilhamtubagus/urlShortener/utils"
	"github.com/labstack/echo/v4"
)

type AuthenticationService interface {
	SignIn(user *entity.User) (*entity.Token, error)
	GoogleSignIn(credential string) (*entity.Token, error)
	RefreshToken(refreshToken string) (*entity.Token, error)
}
type authenticationService struct {
	userService         UserService
	oauth2GoogleService Oauth2GoogleService
	hash                utils.Hash
}

func (a authenticationService) SignIn(user *entity.User) (*entity.Token, error) {
	searchedUser, err := a.userService.FindUserByEmail(user.Email)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, dto.NewDefaultResponse("unexpected database error", http.StatusInternalServerError))
	}
	const EMPTY_STRING = ""
	if searchedUser == nil || searchedUser.Password == EMPTY_STRING {
		return nil, echo.NewHTTPError(http.StatusNotFound, dto.NewDefaultResponse("user was not found", http.StatusNotFound))
	}
	if err := a.hash.CompareHash(user.Password, searchedUser.Password); err != nil {
		fmt.Println(err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, dto.NewDefaultResponse("password does not match", http.StatusBadRequest))
	}
	claims, err := searchedUser.CreateClaims()
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, dto.NewDefaultResponse("unexpected server error", http.StatusInternalServerError))
	}
	token, err := claims.GenerateJwt()
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, dto.NewDefaultResponse("unexpected server error", http.StatusInternalServerError))
	}
	return token, nil
}

func (a authenticationService) GoogleSignIn(credential string) (*entity.Token, error) {
	googleClaims, err := a.oauth2GoogleService.VerifyToken(credential)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	user, err := a.userService.FindUserByEmail(googleClaims.Email)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, dto.NewDefaultResponse("unexpected database error", http.StatusInternalServerError))
	}
	if user == nil {
		//insert new user into database
		user = &entity.User{Name: googleClaims.Name, Email: googleClaims.Email, Sub: googleClaims.Sub, Status: constant.ACTIVE, Role: constant.MEMBER}
		err := a.userService.SaveUser(user)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError, dto.NewDefaultResponse("unexpected database error", http.StatusInternalServerError))
		}
	}
	claims, err := user.CreateClaims()
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, dto.NewDefaultResponse("unexpected server error", http.StatusInternalServerError))
	}

	token, err := claims.GenerateJwt()
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, dto.NewDefaultResponse("unexpected server error", http.StatusInternalServerError))
	}
	return token, nil
}

func (a authenticationService) RefreshToken(refreshToken string) (*entity.Token, error) {
	parsedRefreshToken, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		secret := os.Getenv("TOKEN_SECRET")
		if secret == "" {
			return nil, errors.New("token secret has not been set")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := parsedRefreshToken.Claims.(jwt.MapClaims)
	if !ok || !parsedRefreshToken.Valid {
		return nil, errors.New("Refresh token invalid")
	}
	email := fmt.Sprintf("%v", claims["sub"])
	searchedUser, err := a.userService.FindUserByEmail(email)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, dto.NewDefaultResponse("unexpected database error", http.StatusInternalServerError))
	}
	if searchedUser == nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, dto.NewDefaultResponse("user was not found", http.StatusNotFound))
	}
	refreshedClaims, err := searchedUser.CreateClaims()
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, dto.NewDefaultResponse("unexpected server error", http.StatusInternalServerError))
	}

	token, err := refreshedClaims.GenerateJwt()
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, dto.NewDefaultResponse("unexpected server error", http.StatusInternalServerError))
	}
	return token, nil
}

func NewAuthenticationService(userService UserService, hash utils.Hash, oauth2GoogleService Oauth2GoogleService) AuthenticationService {
	return authenticationService{userService: userService, hash: hash, oauth2GoogleService: oauth2GoogleService}
}
