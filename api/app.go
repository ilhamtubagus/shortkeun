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
	"fmt"
	"os"

	openAPIMiddleware "github.com/go-openapi/runtime/middleware"
	"github.com/ilhamtubagus/urlShortener/authentication"
	"github.com/ilhamtubagus/urlShortener/lib"
	"github.com/ilhamtubagus/urlShortener/user"
	"github.com/kamva/mgm/v3"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func StartApp(e *echo.Echo) {
	//add CustomValidator into echo context
	e.Validator = lib.NewCustomValidator()

	//hash library instantiation
	bcryptHash := lib.NewBcryptHash()

	//collections instantiation based on entities
	userCollection := mgm.Coll(new(user.User))

	//repositories instantiation
	userRepository := user.NewUserRepository(userCollection)

	//service instantiation
	userService := user.NewUserService(userRepository, bcryptHash)
	authenticationService := authentication.NewAuthenticationService(userService, bcryptHash)

	//handlers instantiation
	authenticationController := authentication.NewAuthenticationController(authenticationService)
	userController := user.NewUserController(userService)

	//routes definition
	e.GET("", func(c echo.Context) error {
		return c.File("./public/demoSignInWithGoogle.html")
	})
	//serve API Docs
	e.GET("/swagger.yaml", func(c echo.Context) error {
		return c.File("swagger.yaml")
	})
	opts := openAPIMiddleware.RedocOpts{SpecURL: "/swagger.yaml"}
	redocMiddleware := openAPIMiddleware.Redoc(opts, nil)
	eRedocMiddleware := echo.WrapHandler(redocMiddleware)
	e.GET("/docs", eRedocMiddleware)
	// authentication routes
	e.POST("/tokens", authenticationController.SignIn)
	e.POST("/tokens/google", authenticationController.GoogleSignIn)
	e.PATCH("/users/activation-code", userController.RequestActivationCode)
	e.POST("/users", userController.Register)

	restrictedRoutes := e.Group("/users/status")
	secret := os.Getenv("TOKEN_SECRET")
	if secret == "" {
		fmt.Fprintf(os.Stderr, "error: %s\n", "Token secret not found")
		os.Exit(1)
	}
	jwtConfig := echoMiddleware.JWTConfig{
		Claims:      &authentication.Claims{},
		SigningKey:  []byte(secret),
		TokenLookup: "header:" + echo.HeaderAuthorization,
		AuthScheme:  "Bearer",
	}
	restrictedRoutes.Use(echoMiddleware.JWTWithConfig(jwtConfig))
	restrictedRoutes.PATCH("", userController.ActivateAccount)
}
