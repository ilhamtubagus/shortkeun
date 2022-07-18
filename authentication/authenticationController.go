package authentication

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/ilhamtubagus/urlShortener/authentication/dto"
	commonDto "github.com/ilhamtubagus/urlShortener/dto"
	"github.com/ilhamtubagus/urlShortener/email"
	"github.com/ilhamtubagus/urlShortener/lib"
	"github.com/labstack/echo/v4"
)

type AuthenticationController struct {
	authenticationService AuthenticationService
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
func (controller AuthenticationController) Register(c echo.Context) error {
	registrant := new(dto.RegistrationRequestBody)
	if err := c.Bind(&registrant); err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, commonDto.NewDefaultResponse("failed to parse request body", http.StatusBadRequest))
	}

	if err := c.Validate(registrant); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity,
			commonDto.NewValidationError("validation failed", lib.MapError(err), http.StatusUnprocessableEntity))
	}
	user, err := controller.authenticationService.Register(registrant.ConvertToEntity())
	if err != nil {
		return err
	}
	//send email registration with activation code
	ipE := echo.ExtractIPDirect()
	now := time.Now()
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
	return c.JSON(http.StatusCreated, &user)
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
func (controller AuthenticationController) SignIn(c echo.Context) error {
	var credential dto.SignInRequestDefaultBody
	err := c.Bind(&credential)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, commonDto.NewDefaultResponse("failed to parse request body", http.StatusBadRequest))
	}
	//dto validation
	if err := c.Validate(credential); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity,
			commonDto.NewValidationError("validation failed", lib.MapError(err), http.StatusUnprocessableEntity))
	}
	token, err := controller.authenticationService.SignIn(credential.ConvertToEntity())
	if err != nil {
		return err
	}
	tokenResponse := dto.TokenResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
	return c.JSON(200, &dto.SignInResponseBody{Message: "signin succeeded", Token: tokenResponse})
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
func (controller AuthenticationController) GoogleSignIn(c echo.Context) error {
	var credential dto.GoogleSignInRequestBody
	err := c.Bind(&credential)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, commonDto.NewDefaultResponse("failed to parse request body", http.StatusBadRequest))
	}
	//dto validation
	if err := c.Validate(credential); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity,
			commonDto.NewValidationError("validation failed", lib.MapError(err), http.StatusUnprocessableEntity))
	}
	// decode and verify id token credential
	googleTokenInfo, err := VerifyToken(credential.Credential)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	token, err := controller.authenticationService.GoogleSignIn(googleTokenInfo)
	if err != nil {
		return err
	}
	tokenResponse := dto.TokenResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
	return c.JSON(200, &dto.SignInResponseBody{Message: "signin succeeded", Token: tokenResponse})
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
func (controller AuthenticationController) RequestActivationCode(c echo.Context) error {
	newActivationCodeRequest := new(dto.NewActivationCodeRequestBody)
	if err := c.Bind(&newActivationCodeRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, commonDto.NewDefaultResponse("failed to parse request body", http.StatusBadRequest))
	}
	// dto validation
	if err := c.Validate(newActivationCodeRequest); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity,
			commonDto.NewValidationError("validation failed", lib.MapError(err), http.StatusUnprocessableEntity))
	}

	user := newActivationCodeRequest.ConvertToEntity()
	user, err := controller.authenticationService.RequestActivationCode(user)
	if err != nil {
		return err
	}
	//send email registration with activation code
	ipE := echo.ExtractIPDirect()
	now := time.Now()
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
	return c.JSON(http.StatusCreated, commonDto.NewDefaultResponse("activation code sent", http.StatusOK))
}

func NewAuthenticationController(authenticationService AuthenticationService) AuthenticationController {
	return AuthenticationController{
		authenticationService: authenticationService,
	}
}
