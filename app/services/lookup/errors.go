package lookup

import "errors"

var (
	// ErrInvalidCountryID is error given when a country ID is not valid
	ErrInvalidCountryID = errors.New("Invalid Country ID")

	// ErrCountryNotFound is error given when country is not found
	ErrCountryNotFound = errors.New("No country was found")

	// ErrInvalidStateID is error given when a state ID is not valid
	ErrInvalidStateID = errors.New("Invalid State ID")

	// ErrStateNotFound is error given when a state is not found
	ErrStateNotFound = errors.New("No state was found")
)
