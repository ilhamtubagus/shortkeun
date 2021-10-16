// Package classification Shorkeun API.
// Documentation for Shorkeun API.
// Why do we need to shorten the URL? Is it something necessary?
// Well, there are many advantages that shortening of URL provides.
// A very basic advantage could be that users tend to make very few mistakes while copying the URL if it is not shortened.
// Secondly, they surely save a lot of space when used or printed.
// Moreover, it could be used if someone wishes not to use the original URL or want to hide the original one.
// Terms Of Service:
// There are no TOS at this moment, use at your own risk we take no responsibility
// 	Schemes: http, https
//     Host: localhost
//     BasePath: /
//     Version: 1.0.0
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: Ilham Tubagus Arfian<ilhamta27@gmail.com> https://github.com/ilhamtubagus
//     Consumes:
//     - application/json
// 		Produces:
//     - application/json
//     SecurityDefinitions:
//     Bearer-Token:
//          type: apiKey
//          name: Authorization
//          in: header
// swagger:meta
package api

import (
	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/ilhamtubagus/urlShortener/api/handlers"
	"github.com/ilhamtubagus/urlShortener/entities"
	"github.com/ilhamtubagus/urlShortener/lib"
	"github.com/ilhamtubagus/urlShortener/repository"
	"github.com/ilhamtubagus/urlShortener/service"
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

	//libs instantiation
	bcrypHasher := lib.NewBcryptHasher()
	//services instantiation
	authService := service.NewAuthService(userRepository, bcrypHasher)
	//handlers instantiation
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userRepository)

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
	// auth routes
	e.POST("/auth/signin", authHandler.SignIn)
	e.POST("/auth/signin/google", authHandler.GoogleSignIn)
	e.PATCH("/auth/signin/activation-code", authHandler.RequestActivationCode)
	e.POST("/auth/register", authHandler.Register)
	// user routes
	e.PATCH("/user/status", userHandler.ActivateAccount)
}
