package api

import (
	"github.com/ilhamtubagus/urlShortener/api/handlers"
	"github.com/ilhamtubagus/urlShortener/lib"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

func StartApp(e *echo.Echo, dbClient *mongo.Client) {
	//add CustomValidator
	e.Validator = lib.NewCustomValidator()

	//handlers instantiation
	authHandler := handlers.AuthHandler{}

	//routes definition
	e.GET("/", func(c echo.Context) error {
		return c.File("./public/index.html")
	})
	e.POST("/auth/login/google", authHandler.GoogleSignIn)
}
