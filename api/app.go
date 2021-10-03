package api

import (
	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/ilhamtubagus/urlShortener/api/handlers"
	"github.com/ilhamtubagus/urlShortener/entities"
	"github.com/ilhamtubagus/urlShortener/lib"
	"github.com/ilhamtubagus/urlShortener/repository"
	"github.com/kamva/mgm/v3"
	"github.com/labstack/echo/v4"
)

func StartApp(e *echo.Echo) {
	//add CustomValidator into echo context
	e.Validator = lib.NewCustomValidator()

	//collections instantiation based on entities
	userCollection := mgm.Coll(new(entities.User))

	//repositories instantiation
	userRepository := repository.NewUserRepository(userCollection)

	//handlers instantiation
	authHandler := handlers.NewAuthHandler(userRepository)

	//routes definition
	e.GET("", func(c echo.Context) error {
		return c.File("./public/demoSignInWithGoogle.html")
	})
	//serve API Docs
	e.GET("/swagger.yaml", func(c echo.Context) error {
		return c.File("swagger.yaml")
	})
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	redocMiddleware := middleware.Redoc(opts, nil)
	eRedocMiddleware := echo.WrapHandler(redocMiddleware)
	e.GET("/docs", eRedocMiddleware)
	e.POST("/auth/signin", authHandler.SignIn)
	e.POST("/auth/signin/google", authHandler.GoogleSignIn)
	e.PATCH("/auth/signin/activation-code", authHandler.RequestActivationCode)
	e.POST("/auth/register", authHandler.Register)
}
