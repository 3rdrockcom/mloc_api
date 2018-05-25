package auth

import (
	API "github.com/epointpayment/mloc_api_go/app/services/api"
	Customer "github.com/epointpayment/mloc_api_go/app/services/customer"

	"github.com/labstack/echo"
)

func DefaultValidator(key string, customerUniqueIDFieldName string, c echo.Context) (isValid bool, err error) {
	// Initialize API service
	sa := API.New()

	// Get customer key
	entry, err := sa.GetCustomerKey(key)
	if err != nil {
		return
	}

	// Check if key is a match
	if key != entry.Key {
		err = API.ErrInvalidAPIKey
		return
	}

	// If customer unique ID field name is given, allow access to non-user related (lookup) requests
	if customerUniqueIDFieldName == "" {
		// Is for non-user related requests authorized
		isValid = true
		return
	}

	// Get customer unique ID
	var customerUniqueID string
	method := c.Request().Method
	switch method {
	case "GET":
		customerUniqueID = c.QueryParam(customerUniqueIDFieldName)
	case "POST":
		customerUniqueID = c.FormValue(customerUniqueIDFieldName)
	}

	// Initialize customer service
	sc, err := Customer.New(*entry.CustomerID)
	if err != nil {
		return
	}

	// Get customer information
	customer, err := sc.Info().Get()
	if err != nil {
		return
	}

	// Check if customer unique ID is a match
	if customer.CustomerUniqueID != customerUniqueID {
		err = Customer.ErrInvalidUniqueCustomerID
		return
	}

	// Pass user information to context
	c.Set("customerID", *entry.CustomerID)

	// User is authorized
	isValid = true

	return
}

func RegistrationValidator(key string, customerUniqueIDFieldName string, c echo.Context) (isValid bool, err error) {
	// Initialize API service
	sa := API.New()

	// Get API key for registration
	entry, err := sa.GetRegistrationKey()
	if err != nil {
		return
	}

	// Check if key is a match
	if key != entry.Key {
		err = API.ErrInvalidAPIKey
		return
	}

	// User is authorized
	isValid = true
	return
}

// LoginValidator is a validator used for key auth middleware
func LoginValidator(key string, customerUniqueIDFieldName string, c echo.Context) (isValid bool, err error) {
	// Initialize API service
	sa := API.New()

	// Get API key for login
	entry, err := sa.GetLoginKey()
	if err != nil {
		return
	}

	// Check if key is a match
	if key != entry.Key {
		err = API.ErrInvalidAPIKey
		return
	}

	// User is authorized
	isValid = true
	return
}
