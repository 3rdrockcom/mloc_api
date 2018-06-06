package router

import (
	"net/http"

	"github.com/epointpayment/mloc_api_go/app/controllers"
	API "github.com/epointpayment/mloc_api_go/app/services/api"
	Customer "github.com/epointpayment/mloc_api_go/app/services/customer"
	Lookup "github.com/epointpayment/mloc_api_go/app/services/lookup"

	"github.com/labstack/echo"
)

// appendErrorHandler handles errors for the router
func (r *Router) appendErrorHandler() {
	r.e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		message := err.Error()
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			message = he.Message.(string)
		}

		// Override status code based on error responses
		switch message {
		// API Service
		case API.ErrInvalidAPIKey.Error():
			code = http.StatusForbidden

		// Customer Service
		case Customer.ErrInvalidUniqueCustomerID.Error():
			code = http.StatusForbidden
		case Customer.ErrCustomerNotFound.Error():
			code = http.StatusNotFound

		// Lookup Service
		case Lookup.ErrCountryNotFound.Error():
			code = http.StatusNotFound
		}

		// Send error in a specific format
		controllers.SendErrorResponse(c, code, message)

		// Log errors
		c.Logger().Error(err)
	}
}
