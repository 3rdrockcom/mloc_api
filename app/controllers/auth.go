package controllers

import (
	"net/http"
	"strconv"

	API "github.com/epointpayment/mloc_api_go/app/services/api"

	"github.com/labstack/echo"
)

// GetCustomerKey retrieves an API key from the database
func (co Controllers) GetCustomerKey(c echo.Context) (err error) {
	// Get program ID
	programID := 1

	// Get program customer ID
	programCustomerID, err := strconv.Atoi(c.QueryParam("customer_id"))
	if err != nil {
		err = API.ErrInvalidProgramCustomerID
		return
	}

	// Get program customer mobile
	programCustomerMobile := c.QueryParam("mobile")

	// Get API key
	api := API.New()
	customerAccessKey, err := api.GetCustomerAccessKey(programID, programCustomerID, programCustomerMobile)
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// Send results
	return SendResponse(c, http.StatusOK, customerAccessKey)

}

// GenerateCustomerKey creates a new customer and API key and stores it in the database
func (co Controllers) GenerateCustomerKey(c echo.Context) (err error) {
	// Get program ID
	programID := 1

	// Get program customer ID
	programCustomerID, err := strconv.Atoi(c.QueryParam("customer_id"))
	if err != nil {
		err = API.ErrInvalidProgramCustomerID
		return
	}

	// Get program customer mobile
	programCustomerMobile := c.QueryParam("mobile")

	// Get API key
	api := API.New()
	customerAccessKey, err := api.GenerateCustomerAccessKey(programID, programCustomerID, programCustomerMobile)
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// Send results
	return SendResponse(c, http.StatusOK, customerAccessKey)
}
