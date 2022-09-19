package handler

import (
	"net/http"

	"github.com/ilhamtubagus/urlShortener/domain/entity"
	"github.com/ilhamtubagus/urlShortener/domain/service"
	"github.com/ilhamtubagus/urlShortener/interface/dto"
	"github.com/ilhamtubagus/urlShortener/utils"
	"github.com/labstack/echo/v4"
)

type AuthenticationHandler struct {
	authenticationService service.AuthenticationService
}

// swagger:route POST /tokens auth signIn
//
//		Sign in (default)
//
//		Sign in with email and password
//
//	 Security: None
//		Consumes:
//		 application/json
//		Produces:
//		 application/json
//		Responses:
//		- 422: validationError
//		- 200: signInResponse
//		- 404: defaultResponse
//		- 400: defaultResponse
//		- 500: defaultResponse
func (handler AuthenticationHandler) SignIn(c echo.Context) error {
	var credential dto.SignInRequestDefaultBody
	err := c.Bind(&credential)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, dto.NewDefaultResponse("failed to parse request body", http.StatusBadRequest))
	}
	//dto validation
	if err := c.Validate(credential); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity,
			dto.NewValidationError("validation failed", utils.MapError(err), http.StatusUnprocessableEntity))
	}
	token, err := handler.authenticationService.SignIn(&entity.User{
		Email:    credential.Email,
		Password: credential.Password,
	})
	if err != nil {
		return err
	}
	tokenResponse := dto.TokenResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
	return c.JSON(200, &dto.SignInResponseBody{Message: "Sign in succeeded", Token: tokenResponse})
}

// swagger:route POST /tokens/google auth googleSignIn
//
// # Sign in with google account
//
// Sign in with google account.
// If user has not registered then registration process will be performed.
//
// Consumes:
// - application/json
// Produces:
// - application/json
// Responses:
// - 422: validationError
// - 200: signInResponse
// - 404: defaultResponse
// - 400: defaultResponse
// - 500: defaultResponse
func (handler AuthenticationHandler) GoogleSignIn(c echo.Context) error {
	var credential dto.GoogleSignInRequestBody
	err := c.Bind(&credential)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, dto.NewDefaultResponse("failed to parse request body", http.StatusBadRequest))
	}
	//dto validation
	if err := c.Validate(credential); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity,
			dto.NewValidationError("validation failed", utils.MapError(err), http.StatusUnprocessableEntity))
	}
	// decode and verify id token credential
	token, err := handler.authenticationService.GoogleSignIn(credential.Credential)
	if err != nil {
		return err
	}
	tokenResponse := dto.TokenResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
	return c.JSON(200, &dto.SignInResponseBody{Message: "Sign in succeeded", Token: tokenResponse})
}

func (handler AuthenticationHandler) RefreshToken(c echo.Context) error {
	var credential dto.RefreshTokenRequestBody
	err := c.Bind(&credential)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, dto.NewDefaultResponse("failed to parse request body", http.StatusBadRequest))
	}
	//dto validation
	if err := c.Validate(credential); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity,
			dto.NewValidationError("validation failed", utils.MapError(err), http.StatusUnprocessableEntity))
	}
	pairToken, err := handler.authenticationService.RefreshToken(credential.RefreshToken)
	if err != nil {
		return err
	}
	tokenResponse := dto.TokenResponse{
		AccessToken:  pairToken.AccessToken,
		RefreshToken: pairToken.RefreshToken,
	}
	return c.JSON(200, &dto.SignInResponseBody{Message: "Token refreshed", Token: tokenResponse})
}

func NewAuthenticationHandler(authenticationService service.AuthenticationService) AuthenticationHandler {
	return AuthenticationHandler{
		authenticationService: authenticationService,
	}
}
