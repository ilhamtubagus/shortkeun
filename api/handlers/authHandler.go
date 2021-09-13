package handlers

import (
	"github.com/ilhamtubagus/urlShortener/dto"
	"github.com/ilhamtubagus/urlShortener/entities"
	"github.com/ilhamtubagus/urlShortener/lib"
	"github.com/ilhamtubagus/urlShortener/repositories"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	userRepository repositories.UserRepositories
}

func (a AuthHandler) GoogleSignIn(c echo.Context) error {
	//dto binding
	code := new(dto.SignInRequestGoogle)
	if err := c.Bind(&code); err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
	}
	//dto validation
	if err := c.Validate(code); err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, &dto.ValidationErrorResponse{Message: "Bad Request", Errors: lib.MapError(err)})
	}

	//todo check if signed in user has been registered (present in db)
	return c.JSON(200, code)
}

func (a AuthHandler) DefaultSignIn(c echo.Context) error {
	return c.JSON(200, c.Path())
}

func (a AuthHandler) Register(c echo.Context) error {
	//dto binding
	registrant := new(dto.RegisterRequest)
	if err := c.Bind(&registrant); err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
	}
	//dto validation
	if err := c.Validate(registrant); err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code,
			&dto.ValidationErrorResponse{
				Message: "Bad Request",
				Errors:  lib.MapError(err)})
	}
	/* domain validation
	*  1. email must be unique for each users
	 */

	if user, err := a.userRepository.FindUserByEmail(registrant.Email); err != nil {
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, err.Error())
	} else if user != nil {
		return echo.NewHTTPError(echo.ErrForbidden.Code,
			&dto.ValidationErrorResponse{
				Message: "Bad Request",
				Errors: &[]lib.ValidationError{
					{Field: "email", Message: "email has been registered"}}})
	}

	//hash password
	hasher := lib.NewBcryptHasher()
	hashedPassword, err := hasher.MakeHash(registrant.Password)
	if err != nil {
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, err.Error())
	}
	//create user struct
	user := &entities.User{Name: registrant.Name, Email: registrant.Email, Password: hashedPassword}
	//save user in repository
	inserted, err := a.userRepository.AddNewUser(user)
	if err != nil {
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, err.Error())
	}
	//generate activation code
	//send email registration
	return c.JSON(200, inserted)
}
func NewAuthHandler(userRepository repositories.UserRepositories) AuthHandler {
	return AuthHandler{userRepository: userRepository}
}
