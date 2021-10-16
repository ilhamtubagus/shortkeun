package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/ilhamtubagus/urlShortener/dto"
	"github.com/ilhamtubagus/urlShortener/email"
	"github.com/ilhamtubagus/urlShortener/lib"
	"github.com/ilhamtubagus/urlShortener/service"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService *service.AuthService
}

// swagger:route POST /auth/signin auth signIn
//
//	Sign in (default)
//
//	Sign in with email and password
//
//  Security: None
//	Consumes:
// 	 application/json
// 	Produces:
// 	 application/json
// 	Responses:
// 	- 422: validationError
//	- 200: signInResponse
//	- 404: defaultResponse
//	- 400: defaultResponse
//	- 500: defaultResponse
func (a AuthHandler) SignIn(c echo.Context) error {
	var credential dto.SignInRequestDefaultBody
	err := c.Bind(&credential)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, dto.NewDefaultResponse("failed to parse request body"))
	}
	//dto validation
	if err := c.Validate(credential); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity,
			dto.NewValidationError("validation failed", lib.MapError(err)))
	}
	token, appErr := a.authService.SignIn(&credential)
	if appErr != nil {
		return echo.NewHTTPError(appErr.StatusCode, dto.NewDefaultResponse(appErr.Err.Error()))
	}
	return c.JSON(http.StatusCreated, &dto.SignInResponseBody{Message: "signin succeeded", Token: *token})
}

//	swagger:route POST /auth/signin/google auth googleSignIn
//
//	Sign in with google account
//
//	Sign in with google account.
//	If user has not registered then registration process will be performed.
//
//	Consumes:
// 	- application/json
// 	Produces:
// 	- application/json
// 	Responses:
// 	- 422: validationError
//	- 200: signInResponse
//	- 404: defaultResponse
//	- 400: defaultResponse
//	- 500: defaultResponse
func (a AuthHandler) GoogleSignIn(c echo.Context) error {
	var credential dto.GoogleSignInRequestBody
	err := c.Bind(&credential)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, dto.NewDefaultResponse("failed to parse request body"))
	}
	//dto validation
	if err := c.Validate(credential); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity,
			dto.NewValidationError("validation failed", lib.MapError(err)))
	}
	token, appErr := a.authService.GoogleSignIn(&credential)
	if appErr != nil {
		return echo.NewHTTPError(appErr.StatusCode, dto.NewDefaultResponse(appErr.Err.Error()))
	}
	return c.JSON(http.StatusCreated, &dto.SignInResponseBody{Message: "signin succeeded", Token: *token})
}

//	swagger:route POST /auth/signin/register auth register
//
//  Register new account
//
//	Register new account with email and password.
//	User will be given a code for account activation via email after registration has been performed.
//
//	Consumes:
// 	- application/json
// 	Produces:
// 	- application/json
// 	Responses:
// 	- 422: validationError
//	- 200: defaultResponse
//	- 404: defaultResponse
//	- 400: defaultResponse
//	- 500: defaultResponse
func (a AuthHandler) Register(c echo.Context) error {
	//dto binding
	registrant := new(dto.RegistrationRequestBody)
	if err := c.Bind(&registrant); err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, dto.NewDefaultResponse("failed to parse request body"))
	}
	//dto validation
	if err := c.Validate(registrant); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity,
			dto.NewValidationError("validation failed", lib.MapError(err)))
	}
	user, appErr := a.authService.Register(registrant)
	if appErr != nil {
		return echo.NewHTTPError(appErr.StatusCode, appErr.Err.Error())
	}
	//send email registration with activation code
	now := time.Now()
	ipE := echo.ExtractIPDirect()
	pathToTemplate, _ := filepath.Abs("./email/template/registrationMail.html")
	attachment, _ := filepath.Abs("./logo.png")
	emailBody := email.RegistrationMailBody{
		UserAgent: c.Request().UserAgent(),
		IP:        ipE(c.Request()),
		DateTime:  now.Format("Monday, 02-Jan-06 15:04:05 MST"),
		Code:      user.ActivationCode.Code,
		ExpireAt:  user.ActivationCode.ExpireAt.Format("Monday, 02-Jan-06 15:04:05 MST"),
	}
	// asynchronously send email registration
	go func() {
		err := lib.SendHTMLMail([]string{user.Email}, "Activate Your Account", emailBody, pathToTemplate, []string{attachment})
		if err != nil {
			c.Logger().Error(fmt.Sprintf("failed to send email registration to %s", user.Email))
		}
	}()
	return c.JSON(http.StatusCreated, dto.NewDefaultResponse("registration succeeded, please check your email for account activation"))
}

//	swagger:route POST /auth/signin/activation-code auth getActivationCode
//
//	Get new activation code
//
//	Get new activation code for account activation purpose if the previous activation code has been expired.
//
//	Consumes:
// 	- application/json
// 	Produces:
// 	- application/json
// 	Responses:
// 	- 422: validationError
//	- 201: defaultResponse
//	- 404: defaultResponse
//	- 400: defaultResponse
//	- 500: defaultResponse
func (ah AuthHandler) RequestActivationCode(c echo.Context) error {
	requestCodeAct := new(dto.ActivationCodeRequestBody)
	if err := c.Bind(&requestCodeAct); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, dto.NewDefaultResponse("failed to parse request body"))
	}
	// dto validation
	if err := c.Validate(requestCodeAct); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity,
			dto.NewValidationError("validation failed", lib.MapError(err)))
	}
	user, appErr := ah.authService.RequestActivationCode(requestCodeAct)
	if appErr != nil {
		return echo.NewHTTPError(appErr.StatusCode, appErr.Err.Error())
	}
	now := time.Now()
	//send email registration with activation code
	ipE := echo.ExtractIPDirect()
	pathToTemplate, _ := filepath.Abs("./email/template/activationMail.html")
	attachment, _ := filepath.Abs("./logo.png")
	emailBody := email.RegistrationMailBody{
		UserAgent: c.Request().UserAgent(),
		IP:        ipE(c.Request()),
		DateTime:  now.Format("Monday, 02-Jan-06 15:04:05 MST"),
		Code:      user.ActivationCode.Code,
		ExpireAt:  user.ActivationCode.ExpireAt.Format("Monday, 02-Jan-06 15:04:05 MST"),
	}
	// asynchronously send email registration
	go func() {
		err := lib.SendHTMLMail([]string{user.Email}, "Activate Your Account", emailBody, pathToTemplate, []string{attachment})
		if err != nil {
			c.Logger().Error(fmt.Sprintf("failed to send email registration to %s", user.Email))
		}
	}()
	return c.JSON(http.StatusCreated, dto.NewDefaultResponse("activation code sent"))
}

// todo refresh token

func NewAuthHandler(authService *service.AuthService) AuthHandler {
	return AuthHandler{authService: authService}
}
