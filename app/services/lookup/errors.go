package lookup

import "errors"

var (
	// ErrInvalidCountryID is error given when a country ID is not valid
	ErrInvalidCountryID = errors.New("Invalid Country ID")

	// ErrCountryNotFound is error given when country is not found
	ErrCountryNotFound = errors.New("No country was found")
)
