package handlers

import (
	"fmt"
	"net/http"

	"github.com/ilhamtubagus/urlShortener/dto"
	"github.com/ilhamtubagus/urlShortener/lib"
	"github.com/ilhamtubagus/urlShortener/repository"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userRepository repository.UserRepository
}

//	swagger:route PATCH /user/status user accountActivation
//	Activate user account handler
func (uh UserHandler) ActivateAccount(c echo.Context) error {
	accountActivationReq := new(dto.AccountActivationRequest)
	if err := c.Bind(&accountActivationReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, dto.NewDefaultResponse("failed to parse request body"))
	}
	// dto validation
	if err := c.Validate(accountActivationReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			dto.NewValidationError("Bad Request", lib.MapError(err)))
	}
	fmt.Println(accountActivationReq)
	return nil
}
