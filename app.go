package main

import (
	"fmt"
	"net/http"
	"os"

	openAPIMiddleware "github.com/go-openapi/runtime/middleware"
	"github.com/golang-jwt/jwt"
	"github.com/ilhamtubagus/urlShortener/domain/entity"
	"github.com/ilhamtubagus/urlShortener/domain/service"
	"github.com/ilhamtubagus/urlShortener/infrastructure/persistence"
	"github.com/ilhamtubagus/urlShortener/interface/handler"
	"github.com/ilhamtubagus/urlShortener/utils"
	"github.com/kamva/mgm/v3"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func restricted(c echo.Context) error {
	user := c.(*utils.AuthenticatedContext)
	return c.String(http.StatusOK, "Welcome "+user.Role+"!")
}

// Package classification Shorkeun API.
// Documentation for Shorkeun API.
// Why do we need to shorten the URL? Is it something necessary?
// Well, there are many advantages that shortening of URL provides.
// A very basic advantage could be that users tend to make very few mistakes while copying the URL if it is not shortened.
// Secondly, they surely save a lot of space when used or printed.
// Moreover, it could be used if someone wishes not to use the original URL or want to hide the original one.
// Terms Of Service:
// There are no TOS at this moment, use at your own risk we take no responsibility
//
//		Schemes: http, https
//	    Host: localhost
//	    BasePath: /
//	    Version: 1.0.0
//	    License: MIT http://opensource.org/licenses/MIT
//	    Contact: Ilham Tubagus Arfian<ilhamta27@gmail.com> https://github.com/ilhamtubagus
//	    Consumes:
//	    - application/json
//			Produces:
//	    - application/json
//	    SecurityDefinitions:
//	    Bearer-Token:
//	         type: apiKey
//	         name: Authorization
//	         in: header
//
// swagger:meta
func InitializeEchoApp(e *echo.Echo) {
	//add CustomValidator into echo context
	e.Validator = utils.NewCustomValidator()

	//hash utils instantiation
	bcryptHash := utils.NewBcryptHash()

	//collections instantiation based on entities
	userCollection := mgm.Coll(new(entity.User))

	//repositories instantiation
	userRepository := persistence.NewUserRepository(userCollection)

	//service instantiation
	userService := service.NewUserService(userRepository, bcryptHash)
	googleClientId := os.Getenv("G_CLIENT_ID")
	oauth2Service := service.NewOauth2GoogleService(googleClientId)
	authenticationService := service.NewAuthenticationService(userService, bcryptHash, oauth2Service)

	//handlers instantiation
	authenticationHandler := handler.NewAuthenticationHandler(authenticationService)
	userHandler := handler.NewUserHandler(userService)

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
	e.POST("/tokens", authenticationHandler.SignIn)
	e.POST("/tokens/google", authenticationHandler.GoogleSignIn)
	e.PATCH("/tokens", authenticationHandler.RefreshToken)
	e.PATCH("/users/activation-code", userHandler.RequestActivationCode)
	e.POST("/users", userHandler.Register)
	e.PATCH("/users/status", userHandler.ActivateAccount)
	secret := os.Getenv("TOKEN_SECRET")
	if secret == "" {
		fmt.Fprintf(os.Stderr, "error: %s\n", "Token secret not found")
		os.Exit(1)
	}

	// Restricted group
	r := e.Group("/urls")
	jwtConfig := echoMiddleware.JWTConfig{
		Claims:      &entity.Claims{},
		SigningKey:  []byte(secret),
		TokenLookup: "header:" + echo.HeaderAuthorization,
		AuthScheme:  "Bearer",
	}

	r.Use(echoMiddleware.JWTWithConfig(jwtConfig))
	r.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(*entity.Claims)
			cc := &utils.AuthenticatedContext{Context: c, Email: claims.Email, Role: claims.Role, Status: claims.Status}
			return next(cc)
		}
	})
	r.POST("", restricted)
}
