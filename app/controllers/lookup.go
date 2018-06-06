package controllers

import (
	"net/http"
	"strconv"

	Lookup "github.com/epointpayment/mloc_api_go/app/services/lookup"

	"github.com/labstack/echo"
)

// GetCountry gets the information of country from database
func (co Controllers) GetCountry(c echo.Context) (err error) {
	// Get country id
	countryID, err := strconv.Atoi(c.QueryParam("country_id"))
	if err != nil {
		err = Lookup.ErrInvalidCountryID

		// Send response
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
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
