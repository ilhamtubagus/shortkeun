package user

import (
	"net/http"

	commonDto "github.com/ilhamtubagus/urlShortener/dto"
	"github.com/ilhamtubagus/urlShortener/lib"
	"github.com/ilhamtubagus/urlShortener/user/dto"
	"github.com/labstack/echo/v4"
)

type UserController struct {
	userService UserService
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
func (controller UserController) ActivateAccount(c echo.Context) error {
	accountActivationRequest := new(dto.AccountActivationRequestBody)
	if err := c.Bind(&accountActivationRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, commonDto.NewDefaultResponse("failed to parse request body", http.StatusBadRequest))
	}
	// dto validation
	if err := c.Validate(accountActivationRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			commonDto.NewValidationError("Bad Request", lib.MapError(err), http.StatusUnprocessableEntity))
	}
	user, err := controller.userService.ActivateAccount(accountActivationRequest.ActivationCode)
	if err != nil {
		return err
	}
	userResponseDto := user.ConvertToDto()
	return c.JSON(http.StatusOK, &userResponseDto)
}
