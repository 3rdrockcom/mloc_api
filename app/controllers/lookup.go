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

// GetCity gets the information about a city
func (co *Controllers) GetCity(c echo.Context) (err error) {
	// Get state code
	stateCode := c.QueryParam("state_code")
	if stateCode == "" {
		err = Lookup.ErrInvalidStateCode
		return
	}

	// Get city id
	queryCityID := c.QueryParam("city_id")
	cityID, err := strconv.Atoi(queryCityID)
	if err != nil && queryCityID != "" {
		err = Lookup.ErrInvalidCityID
		return
	}

	// Initialize lookup service
	ls := Lookup.New()

	// Get city
	if cityID > 0 {
		city, err := ls.GetCity(stateCode, cityID)
		if err != nil {
			return err
		}

		// Send response
		return SendResponse(c, http.StatusOK, city)
	}

	// Get list of cities in a state
	city, err := ls.GetCities(stateCode)
	if err != nil {
		return err
	}

	// Send response
	return SendResponse(c, http.StatusOK, city)
}

// GetIncomeSource gets the information about an income source
func (co *Controllers) GetIncomeSource(c echo.Context) (err error) {
	// Get id
	queryID := c.QueryParam("id")
	id, err := strconv.Atoi(queryID)
	if err != nil && queryID != "" {
		err = Lookup.ErrInvalidIncomeSourceID
		return
	}

	// Initialize lookup service
	ls := Lookup.New()

	// Get income source
	if id > 0 {
		incomeSource, err := ls.GetIncomeSource(id)
		if err != nil {
			return err
		}

		// Send response
		return SendResponse(c, http.StatusOK, incomeSource)
	}

	// Get list of income sources
	incomeSources, err := ls.GetIncomeSources()
	if err != nil {
		return err
	}

	// Send response
	return SendResponse(c, http.StatusOK, incomeSources)
}
