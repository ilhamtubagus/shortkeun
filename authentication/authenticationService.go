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
	usr "github.com/ilhamtubagus/urlShortener/user"
	"github.com/labstack/echo/v4"
)

type AuthenticationService interface {
	Register(user *usr.User) (*usr.User, error)
	SignIn(user *usr.User) (*Token, error)
	RequestActivationCode(user *usr.User) (*usr.User, error)
	GoogleSignIn(googleClaims *GoogleClaims) (*Token, error)
}

type authenticationService struct {
	userService usr.UserService
	hash        lib.Hash
}

func (service authenticationService) Register(user *usr.User) (*usr.User, error) {
	// domain validation
	// email must be unique for each users
	if user, err := service.userService.FindUserByEmail(user.Email); err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	} else if user != nil {
		return nil, echo.NewHTTPError(http.StatusUnprocessableEntity,
			dto.NewValidationError("validation failed", &[]lib.ValidationError{
				{Field: "email", Message: "email has been registered"}}, http.StatusBadRequest))
	}
	// perform password hashing
	hashedPassword, err := service.hash.MakeHash(user.Password)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	user.Password = *hashedPassword
	//generate activation code
	activationCodeExpiry, _ := strconv.Atoi(os.Getenv("ACTIVATION_CODE_EXPIRY_IN_MINUTES"))
	activationCode := lib.RandString(5)
	now := time.Now()
	user.ActivationCode = &usr.ActivationCode{
		Code:     activationCode,
		IssuedAt: now,
		ExpireAt: now.Add(time.Minute * time.Duration(activationCodeExpiry)),
	}
	err = service.userService.SaveUser(user)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return user, nil
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

func (authenticationService authenticationService) RequestActivationCode(user *usr.User) (*usr.User, error) {
	user, err := authenticationService.userService.FindUserByEmail(user.Email)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if user == nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, dto.NewDefaultResponse("user with this email address was not found", http.StatusNotFound))
	}
	if user.Status == usr.SUSPENDED {
		return nil, echo.NewHTTPError(http.StatusBadRequest, dto.NewDefaultResponse("user with this email address is suspended, please contact administrator for further information", http.StatusUnprocessableEntity))
	}
	if user.Status == usr.ACTIVE {
		return nil, echo.NewHTTPError(http.StatusBadRequest, dto.NewDefaultResponse("user with this email address has already been activated", http.StatusBadRequest))
	}
	if user.ActivationCode != nil {
		if user.ActivationCode.ExpireAt.In(time.Now().Location()).After(time.Now()) {
			return nil, echo.NewHTTPError(http.StatusBadRequest, dto.NewDefaultResponse("the previous activation code has not been expired", http.StatusBadRequest))
		}
	}
	// issue new activation code
	now := time.Now()
	activationCode := lib.RandString(5)
	user.ActivationCode = &usr.ActivationCode{
		Code:     activationCode,
		IssuedAt: now,
		ExpireAt: now.Add(5 * time.Minute),
	}
	err = authenticationService.userService.UpdateUser(user)
	if err != nil {
		fmt.Println(err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return user, nil
}
