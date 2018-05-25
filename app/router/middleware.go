package router

import (
	"github.com/epointpayment/mloc_api_go/app/router/middleware/auth"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// appendMiddleware registers middleware
func (r *Router) appendMiddleware() {
	r.e.Use(middleware.Gzip())
	r.e.Use(middleware.Logger())
	r.e.Use(middleware.Recover())
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
