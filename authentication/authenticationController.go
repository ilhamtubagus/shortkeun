package authentication

import (
	"net/http"

	"github.com/ilhamtubagus/urlShortener/authentication/dto"
	commonDto "github.com/ilhamtubagus/urlShortener/dto"
	"github.com/ilhamtubagus/urlShortener/lib"
	"github.com/ilhamtubagus/urlShortener/user"
	"github.com/labstack/echo/v4"
)

type authenticationController struct {
	userService user.UserService
}

func (controller authenticationController) Register(c echo.Context) error {
	registrant := new(dto.RegistrationRequestBody)
	if err := c.Bind(&registrant); err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, commonDto.NewDefaultResponse("failed to parse request body"))
	}

	if err := c.Validate(registrant); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity,
			commonDto.NewValidationError("validation failed", lib.MapError(err)))
	}

	
}

func NewAuthenticationController(userService user.UserService) authenticationController {
	return authenticationController{
		userService: userService,
	}
}
