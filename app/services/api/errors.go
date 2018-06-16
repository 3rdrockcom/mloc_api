package api

import "errors"

var (
	// ErrInvalidAPIKey is an error given when the requestor uses an invalid API key
	ErrInvalidAPIKey = errors.New("Invalid API Key")

	// ErrInvalidProgramCustomerID is an error shown when customer ID is not a valid
	ErrInvalidProgramCustomerID = errors.New("Invalid Customer ID")

	// ErrCustomerExists is an error given when the customer was already created
	ErrCustomerExists = errors.New("Customer already has an existing id and key")
)
