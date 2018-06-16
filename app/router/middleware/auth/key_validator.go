package auth

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"

	API "github.com/epointpayment/mloc_api_go/app/services/api"
	Customer "github.com/epointpayment/mloc_api_go/app/services/customer"

	"github.com/labstack/echo"
)

type Request struct {
	LoanAmount float64 `form:"R2" json:"R2"`
}

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

		switch {
		case strings.HasPrefix(c.Request().Header.Get(echo.HeaderContentType), echo.MIMEApplicationJSON):
			// Get request information
			req := c.Request()

			// Create a temporary buffer
			b := bytes.NewBuffer(make([]byte, 0))

			// TeeReader returns a Reader that writes to b what it reads from req.Body
			reader := io.TeeReader(req.Body, b)

			// Unmarshal json request body to map
			payload := make(map[string]interface{})
			if err := json.NewDecoder(reader).Decode(&payload); err != nil {
				break
			}
			req.Body.Close()

			// NopCloser returns a ReadCloser with a no-op Close method wrapping the provided Reader b
			req.Body = ioutil.NopCloser(b)

			// Get customer unique id from payload
			if v, ok := payload[customerUniqueIDFieldName].(string); ok {
				customerUniqueID = v
			}
		default:
			customerUniqueID = c.FormValue(customerUniqueIDFieldName)
		}
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
