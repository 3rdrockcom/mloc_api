package router

import (
	"log"

	"github.com/epointpayment/mloc_api_go/app/config"
	"github.com/epointpayment/mloc_api_go/app/router/middleware/auth"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// appendMiddleware registers middleware
func (r *Router) appendMiddleware() {
	r.e.Pre(middleware.RemoveTrailingSlash())
	r.e.Use(middleware.Logger())
	r.e.Use(middleware.Recover())

	if config.IsDev() {
		r.e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
			log.Println("Request Body:\n" + string(reqBody))
			log.Println("Response Body:\n" + string(resBody))
		}))
	}
}

// mwBasicAuth handles the basic authentication for a specific route
func (r *Router) mwBasicAuth() echo.MiddlewareFunc {
	return middleware.BasicAuth(auth.BasicValidator)
}

func (r *Router) mwKeyAuth(authType string, customerUniqueIDFieldName string) echo.MiddlewareFunc {
	validator := auth.DefaultValidator

	switch authType {
	case "login":
		validator = auth.LoginValidator
	case "registration":
		validator = auth.RegistrationValidator
	}

	return auth.KeyAuthWithConfig(auth.KeyAuthConfig{
		KeyLookup:                 "header:X-API-KEY",
		AuthScheme:                "",
		Validator:                 validator,
		CustomerUniqueIDFieldName: customerUniqueIDFieldName,
	})
}
