// Package classification Petstore API.
//
// the purpose of this application is to provide an application
// that is using plain go code to define an API
//
// This should demonstrate all the possible comment annotations
// that are available to turn go code into a fully compliant swagger 2.0 spec
//
// Terms Of Service:
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//     Schemes: http, https
//     Host: localhost
//     BasePath: /v2
//     Version: 0.0.1
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: John Doe<john.doe@example.com> http://john.doe.com
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - jwt:
//
//     SecurityDefinitions:
//     jwt:
//          type: bearer token
//          name: Authorization
//          in: header
//     oauth2:
//         type: oauth2
//         authorizationUrl: /oauth2/auth
//         tokenUrl: /oauth2/token
//         in: header
//         scopes:
//           bar: foo
//         flow: accessCode
//
//     Extensions:
//     x-meta-value: value
//     x-meta-array:
//       - value1
//       - value2
//     x-meta-array-obj:
//       - name: obj
//         value: field
//
// swagger:meta
package api

import (
	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/ilhamtubagus/urlShortener/api/handlers"
	"github.com/ilhamtubagus/urlShortener/entities"
	"github.com/ilhamtubagus/urlShortener/lib"
	"github.com/ilhamtubagus/urlShortener/repositories"
	"github.com/kamva/mgm/v3"
	"github.com/labstack/echo/v4"
)

func StartApp(e *echo.Echo) {
	//add CustomValidator into echo context
	e.Validator = lib.NewCustomValidator()

	//collections instantiation based on entities
	userCollection := mgm.Coll(new(entities.User))

	//repositories instantiation
	userRepository := repositories.NewUserRepository(userCollection)

	//handlers instantiation
	authHandler := handlers.NewAuthHandler(userRepository)

	//routes definition
	e.GET("/", func(c echo.Context) error {
		return c.File("./public/index.html")
	})
	//serve API Docs
	e.GET("/swagger.yaml", func(c echo.Context) error {
		return c.File("swagger.yaml")
	})
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	redocMiddleware := middleware.Redoc(opts, nil)
	eRedocMiddleware := echo.WrapHandler(redocMiddleware)
	e.GET("/docs", eRedocMiddleware)

	e.POST("/auth/signin/google", authHandler.GoogleSignIn)
	e.POST("/auth/register", authHandler.Register)
}
