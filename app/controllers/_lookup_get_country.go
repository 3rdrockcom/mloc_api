package controllers

import (
	"net/http"
	"strconv"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/labstack/echo"
)

// Countries is an array of Country entries
type Countries []Country

// Country contains information about a coutryID
type Country struct {
	CountryID        int     `json:"country_id"` // CountryID should be match from database
	Name             string  `json:"name"`
	IsoCode2         string  `db:"iso_code_2" json:"iso_code_2"`
	IsoCode3         string  `db:"iso_code_3" json:"iso_code_3"`
	AddressFormat    string  `json:"address_format"`
	PostCodeRequired int     `json:"postcode_required"`
	Status           int     `json:"status"`
	MobilePrefix     *string `json:"mobile_prefix"`
}

// TableName gets the name of the database table
func (c Country) TableName() string {
	return "tblcountry"
}

var (
	// ErrCountryNotFound is error given when country is not found
	ErrCountryNotFound = "No country was found."
)

// GetCountry gets the information of country from database
func (co *Controllers) GetCountry(c echo.Context) error {
	queryCountryID := c.QueryParam("country_id")

	// if the queryCountryID is empty, then display all information of country table from database
	if len(queryCountryID) == 0 {
		countries := Countries{}

		err := db.Select().
			From("tblcountry").
			All(&countries)

		if err != nil {
			message := ErrCountryNotFound
			return SendErrorResponse(c, http.StatusNotFound, message)
		}

		return SendResponse(c, http.StatusOK, countries)
	}

	// convert queryCountryId to int, if it is not a integer,then return error message
	countryID, err := strconv.Atoi(queryCountryID)
	if err != nil {
		message := ErrCountryNotFound
		return SendErrorResponse(c, http.StatusNotFound, message)
	}

	// display all information of table country from database if it is exist
	country := Country{}
	err = db.Select().
		From("tblcountry").
		Where(dbx.HashExp{"country_id": countryID}).
		One(&country)
	if err != nil {
		message := ErrCountryNotFound
		return SendErrorResponse(c, http.StatusNotFound, message)
	}
	return SendResponse(c, http.StatusOK, country)

}
