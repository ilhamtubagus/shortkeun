package handlers

import (
	"net/http"

	"github.com/ilhamtubagus/urlShortener/dto"
	"github.com/ilhamtubagus/urlShortener/entities"
	"github.com/ilhamtubagus/urlShortener/lib"
	"github.com/ilhamtubagus/urlShortener/service"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService *service.UserService
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
	user := c.Get("user")
	claims, _ := user.(*entities.Claims)
	userId := claims.UserId
	appErr := uh.userService.ActivateAccount(userId, accountActivationReq.ActivationCode)
	if appErr != nil {
		return echo.NewHTTPError(appErr.StatusCode, dto.NewDefaultResponse(appErr.Err.Error()))
	}

	return c.JSON(http.StatusOK, dto.NewDefaultResponse("account activated"))
}

func NewUserHandler(userService *service.UserService) UserHandler {
	return UserHandler{userService: userService}
}
