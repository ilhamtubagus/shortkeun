package authentication

import (
	"net/http"

	"github.com/ilhamtubagus/urlShortener/dto"
	"github.com/labstack/echo/v4"
)

func JWTErrorHandler(err error, c echo.Context) error {
	return c.JSON(http.StatusUnauthorized, dto.DefaultResponseBody{Message: "Invalid token"})
}
