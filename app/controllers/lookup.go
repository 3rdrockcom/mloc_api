package controllers

import (
	"net/http"
	"strconv"

	Lookup "github.com/epointpayment/mloc_api_go/app/services/lookup"

	"github.com/labstack/echo"
)

// GetCountry gets the information about a country
func (co Controllers) GetCountry(c echo.Context) (err error) {
	// Get country id
	queryCountryID := c.QueryParam("country_id")
	countryID, err := strconv.Atoi(queryCountryID)
	if err != nil && queryCountryID != "" {
		err = Lookup.ErrInvalidCountryID
		return
	}

	// Initialize lookup service
	ls := Lookup.New()

	// Get country
	if countryID > 0 {
		country, err := ls.GetCountry(countryID)
		if err != nil {
			return err
		}

		// Send response
		return SendResponse(c, http.StatusOK, country)
	}

	// Get list of countries
	countries, err := ls.GetCountries()
	if err != nil {
		return err
	}

	// Send response
	return SendResponse(c, http.StatusOK, countries)
}

// GetState gets the information about a state
func (co *Controllers) GetState(c echo.Context) (err error) {
	// Get country id
	countryID, err := strconv.Atoi(c.QueryParam("country_id"))
	if err != nil {
		err = Lookup.ErrInvalidCountryID
		return
	}

	// Get state id
	queryStateID := c.QueryParam("state_id")
	stateID, err := strconv.Atoi(queryStateID)
	if err != nil && queryStateID != "" {
		err = Lookup.ErrInvalidStateID
		return
	}

	// Initialize lookup service
	ls := Lookup.New()

	// Get state
	if stateID > 0 {
		state, err := ls.GetState(countryID, stateID)
		if err != nil {
			return err
		}

		// Send response
		return SendResponse(c, http.StatusOK, state)
	}

	// Get list of states for a country
	states, err := ls.GetStates(countryID)
	if err != nil {
		return err
	}

	// Send response
	return SendResponse(c, http.StatusOK, states)
}
