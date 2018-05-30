package controllers

import (
	"net/http"
	"strconv"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/labstack/echo"
)

// Cities is an array of City entries
type Cities []City

// City contains information about a CityID
type City struct {
	CityID    int    `db:"city_id" json:"city_id"`
	City      string `db:"city" json:"city"`
	StateCode string `db:"state_code" json:"state_code"`
}

// TableName gets the name of the database table
func (c City) TableName() string {
	return "tblcity"
}

// IsUpperCaseLetter checks the state_code is upper case letter
func IsUpperCaseLetter(s string) bool {
	for _, r := range s {
		if r < 'A' || r > 'Z' {
			return false
		}
	}
	return true
}

var (
	// ErrCityNotFound is error given when city is not found
	ErrCityNotFound = "No city was found."
)

// GetCity gets the information of city from database
func (co *Controllers) GetCity(c echo.Context) error {

	queryStateCode := c.QueryParam("state_code")
	queryCityID := c.QueryParam("city_id")

	// If the queryStateCode is empty, then it returns error
	if len(queryStateCode) == 0 {
		message := ErrCityNotFound
		return SendErrorResponse(c, http.StatusNotFound, message)
	}

	// check state_code is upper case letter
	uppperCaseLetterQuery := IsUpperCaseLetter(queryStateCode)
	if uppperCaseLetterQuery == false {
		message := ErrCityNotFound
		return SendErrorResponse(c, http.StatusNotFound, message)

	}

	// If queryCityID exists, it displays information in tblcity from database
	// If queryCityID is empty, it dispay all rows base on state code
	if len(queryCityID) == 0 {
		cities := Cities{}

		err := db.Select().
			From("tblcity").
			Where(dbx.HashExp{"state_code": queryStateCode}).
			All(&cities)

		if err != nil {
			message := ErrCityNotFound
			return SendErrorResponse(c, http.StatusNotFound, message)
		}

		return SendResponse(c, http.StatusOK, cities)
	}

	// If queryCityID exists, it displays one row in tblcity from database
	cityID, err := strconv.Atoi(queryCityID)
	if err != nil {
		message := ErrCityNotFound
		return SendErrorResponse(c, http.StatusNotFound, message)
	}

	city := City{}
	err = db.Select().
		From("tblcity").
		Where(dbx.HashExp{"state_code": queryStateCode, "city_id": cityID}).
		One(&city)

	if err != nil {
		message := ErrCityNotFound
		return SendErrorResponse(c, http.StatusNotFound, message)
	}

	return SendResponse(c, http.StatusOK, city)
}
