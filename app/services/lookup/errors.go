package lookup

import "errors"

var (
	// ErrInvalidCountryID is error given when a country ID is not valid
	ErrInvalidCountryID = errors.New("Invalid Country ID")

	// ErrCountryNotFound is error given when country is not found
	ErrCountryNotFound = errors.New("No country was found")

	// ErrInvalidStateID is error given when a state ID is not valid
	ErrInvalidStateID = errors.New("Invalid State ID")

	// ErrInvalidStateCode is error given when a state code is not valid
	ErrInvalidStateCode = errors.New("Invalid State Code")

	// ErrStateNotFound is error given when a state is not found
	ErrStateNotFound = errors.New("No state was found")

	// ErrInvalidCityID is error given when a city ID is not valid
	ErrInvalidCityID = errors.New("Invalid City ID")

	// ErrCityNotFound is error given when a city is not found
	ErrCityNotFound = errors.New("No city was found")
)
