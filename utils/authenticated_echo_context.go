package utils

import "github.com/labstack/echo/v4"

type AuthenticatedContext struct {
	echo.Context
	Email  string
	Role   string
	Status string
}
