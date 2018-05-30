package controllers

import (
	"net/http"
	"strconv"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/labstack/echo"
)

// States is an array of State entries
type States []State

// State contains information about a StateID
type State struct {
	StateID   int    `db:"state_id" json:"state_id"`
	State     string `db:"state" json:"state"`
	StateCode string `db:"state_code" json:"state_code"`
	CountryID *int   `db:"country_id" json:"country_id"`
}

// TableName gets the name of the database table
func (c State) TableName() string {
	return "tblstate"
}

var (
	// ErrStateNotFound is error given when state is not found
	ErrStateNotFound = "No state was found."
)

// GetState gets the information of state in tblstate from database
//returns all rows of the state if the state_id parameter has no value,otherwise single row will be returned
func (co *Controllers) GetState(c echo.Context) error {
	queryCountryID := c.QueryParam("country_id")
	queryStateID := c.QueryParam("state_id")

	// If the queryCountryID is empty, then it displays error
	if len(queryCountryID) == 0 {
		message := ErrStateNotFound
		return SendErrorResponse(c, http.StatusNotFound, message)
	}

	// convert queryCountryID to integer if queryCountryID exists
	countryID, err := strconv.Atoi(queryCountryID)
	if err != nil {
		message := ErrStateNotFound
		return SendErrorResponse(c, http.StatusNotFound, message)
	}

	// If queryStateID is emptry, it dispays all rows base on country ID in tblstate
	if len(queryStateID) == 0 {
		states := States{}

		err := db.Select().
			From("tblstate").
			Where(dbx.HashExp{"country_id": countryID}).
			All(&states)

		if err != nil {
			message := ErrStateNotFound
			return SendErrorResponse(c, http.StatusNotFound, message)
		}

		return SendResponse(c, http.StatusOK, states)
	}

	// If the queryStateID exists, it displays single row of state from database
	stateID, err := strconv.Atoi(queryStateID)
	if err != nil {
		message := ErrStateNotFound
		return SendErrorResponse(c, http.StatusNotFound, message)
	}

	state := State{}
	err = db.Select().
		From("tblstate").
		Where(dbx.HashExp{"state_id": stateID, "country_id": countryID}).
		One(&state)

	if err != nil {
		message := ErrStateNotFound
		return SendErrorResponse(c, http.StatusNotFound, message)
	}

	return SendResponse(c, http.StatusOK, state)
}
