package handlers

import (
	"github.com/ilhamtubagus/urlShortener/dto"
	"github.com/ilhamtubagus/urlShortener/lib"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
}

func (a AuthHandler) GoogleSignIn(c echo.Context) error {
	code := new(dto.SignInRequestGoogle)
	if err := c.Bind(&code); err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
	}

	if err := c.Validate(code); err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, &dto.ValidationErrorResponse{Message: "Bad Request", Errors: lib.MapError(err)})
	}
	return c.JSON(200, code)
}
