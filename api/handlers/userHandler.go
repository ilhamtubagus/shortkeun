package handlers

import (
	"net/http"

	"github.com/ilhamtubagus/urlShortener/dto"
	"github.com/ilhamtubagus/urlShortener/lib"
	"github.com/ilhamtubagus/urlShortener/repository"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userRepository repository.UserRepository
}

// swagger:route PATCH /user/status user accountActivation
//
//	Account activation
//
//	Activate user's account with activation code sent via email. User's status will change to "ACTIVE".
//
//  Security:
// 	- Bearer-Token:
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
func (uh UserHandler) ActivateAccount(c echo.Context) error {
	accountActivationReq := new(dto.AccountActivationRequestBody)
	if err := c.Bind(&accountActivationReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, dto.NewDefaultResponse("failed to parse request body"))
	}
	// dto validation
	if err := c.Validate(accountActivationReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			dto.NewValidationError("Bad Request", lib.MapError(err)))
	}
	// check if activation_code present in users collection
	// check if activation_code is equals with activation_code found in users collection
	return c.JSON(200, accountActivationReq)
}

func NewUserHandler(userRepository repository.UserRepository) UserHandler {
	return UserHandler{userRepository: userRepository}
}
