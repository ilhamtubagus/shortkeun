package authentication

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ilhamtubagus/urlShortener/dto"
	"github.com/ilhamtubagus/urlShortener/lib"
	"github.com/ilhamtubagus/urlShortener/user"
	usr "github.com/ilhamtubagus/urlShortener/user"
	"github.com/labstack/echo/v4"
)

type AuthenticationService interface {
	SignIn(user *usr.User) (*Token, error)
	GoogleSignIn(googleClaims *GoogleClaims) (*Token, error)
}

type authenticationService struct {
	userService usr.UserService
	hash        lib.Hash
}

func (service authenticationService) SignIn(user *usr.User) (*Token, error) {
	searchedUser, err := service.userService.FindUserByEmail(user.Email)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, dto.NewDefaultResponse("unexpected database error", http.StatusInternalServerError))
	}
	const EMPTY_STRING = ""
	if searchedUser == nil || searchedUser.Password == EMPTY_STRING {
		return nil, echo.NewHTTPError(http.StatusNotFound, dto.NewDefaultResponse("user was not found", http.StatusNotFound))
	}
	if err := service.hash.CompareHash(user.Password, searchedUser.Password); err != nil {
		fmt.Println(err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, dto.NewDefaultResponse("password does not match", http.StatusBadRequest))
	}
	hour, _ := strconv.Atoi(os.Getenv("TOKEN_EXP"))
	claims := Claims{
		Role:   searchedUser.Role,
		Email:  searchedUser.Email,
		Status: searchedUser.Status,
		StandardClaims: jwt.StandardClaims{
			//token expires within x hours
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(hour)).Unix(),
			Subject:   searchedUser.ID.String(),
		}}
	token, err := claims.GenerateJwt()
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, dto.NewDefaultResponse("unexpected server error", http.StatusInternalServerError))
	}
	return token, nil
}

func (authenticationService authenticationService) GoogleSignIn(googleClaims *GoogleClaims) (*Token, error) {
	user, err := authenticationService.userService.FindUserByEmail(googleClaims.Email)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, dto.NewDefaultResponse("unexpected database error", http.StatusInternalServerError))
	}
	if user == nil {
		//insert new user into database
		user = &usr.User{Name: googleClaims.Name, Email: googleClaims.Email, Sub: googleClaims.Sub, Status: usr.ACTIVE, Role: usr.MEMBER}
		err := authenticationService.userService.SaveUser(user)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError, dto.NewDefaultResponse("unexpected database error", http.StatusInternalServerError))
		}
	}
	//create our own jwt and send back to client
	hour, _ := strconv.Atoi(os.Getenv("TOKEN_EXP"))
	claims := Claims{
		Role:   user.Role,
		Email:  user.Email,
		Status: user.Status,
		StandardClaims: jwt.StandardClaims{
			//token expires within x hours
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(hour)).Unix(),
			Subject:   user.ID.String(),
		}}
	token, err := claims.GenerateJwt()
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, dto.NewDefaultResponse("unexpected server error", http.StatusInternalServerError))
	}
	return token, nil
}

func NewAuthenticationService(userService user.UserService, hash lib.Hash) AuthenticationService {
	return authenticationService{userService: userService, hash: hash}
}
