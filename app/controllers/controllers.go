package controllers

import (
	"net/http"

	"github.com/epointpayment/mloc_api_go/app/database"
	"github.com/epointpayment/mloc_api_go/app/helpers"
	dbx "github.com/go-ozzo/ozzo-dbx"

	"github.com/labstack/echo"
)

// DB is the database handler
var db *dbx.DB

// Controllers manages the controllers used in the application
type Controllers struct{}

// NewControllers creates an instance of the service
func NewControllers(database *database.Database) *Controllers {
	db = database.GetInstance()
	c := &Controllers{}
	return c
}

// SendResponse sends a response to requestor
func SendResponse(c echo.Context, code int, i interface{}) error {
	return c.JSON(code, i)
}

// SendOKResponse sends a StatusOK (200) response to requestor
func SendOKResponse(c echo.Context, message string) error {
	code := http.StatusOK
	return c.JSON(code, helpers.H{
		"status":        true,
		"response_code": code,
		"message":       message,
	})
}

// SendErrorResponse sends an error response to requestor
func SendErrorResponse(c echo.Context, code int, message string) error {
	return c.JSON(code, helpers.H{
		"status":        false,
		"response_code": code,
		"error":         message,
	})
}
