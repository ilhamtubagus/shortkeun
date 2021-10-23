package middleware

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/ilhamtubagus/urlShortener/entities"
	"github.com/labstack/echo/v4"
)

func IsAuthenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		//extract authorization header
		authorizationHeader := c.Request().Header.Get("Authorization")
		if authorizationHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header is required")
		}
		tokenMatch, _ := regexp.MatchString(`Bearer\s(\S+)`, authorizationHeader)
		if !tokenMatch {
			return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header invalid")
		}
		tokenString := strings.Split(authorizationHeader, " ")[1]
		if tokenString == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header is required")
		}

		//parse into *jwt.Token
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if method, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("signing method invalid")
			} else if method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("signing method invalid")
			}
			return []byte(os.Getenv("TOKEN_SECRET")), nil
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized [failed while authorizing token")
		}
		claims, err := entities.BuildMapClaims(token.Claims)
		if err != nil || !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized [token invalid]")
		}
		c.Set("user", claims)
		return next(c)

	}
}
func IsMember(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user")
		claims, ok := user.(*entities.Claims)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}
		if claims.Role != entities.RoleMember {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized request")
		}
		return next(c)
	}
}

func IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user")
		claims, ok := user.(*entities.Claims)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}
		if claims.Role != entities.RoleAdmin {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized request")
		}
		return next(c)
	}
}
