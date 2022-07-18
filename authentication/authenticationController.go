package authentication

import (
	"net/http"

	"github.com/ilhamtubagus/urlShortener/authentication/dto"
	commonDto "github.com/ilhamtubagus/urlShortener/dto"
	"github.com/ilhamtubagus/urlShortener/lib"
	"github.com/labstack/echo/v4"
)

type AuthenticationController struct {
	authenticationService AuthenticationService
}

// swagger:route POST /tokens auth signIn
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

//	swagger:route POST /tokens/google auth googleSignIn
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

func NewAuthenticationController(authenticationService AuthenticationService) AuthenticationController {
	return AuthenticationController{
		authenticationService: authenticationService,
	}
}
