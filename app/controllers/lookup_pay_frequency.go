package controllers

import (
	"net/http"
	"strconv"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/labstack/echo"
)

// PayFrequencies is an array of PayFrequency entries
type PayFrequencies []PayFrequency

// PayFrequency contains information about a PayFrequencyID
type PayFrequency struct {
	ID          int     `db:"id" json:"id"`
	Description *string `db:"description" json:"description"`
}

// TableName gets the name of the database table
func (c PayFrequency) TableName() string {
	return "tblpayfrequency"
}

var (
	// ErrPayFrequencyNotFound is error given when pay frequency is not found
	ErrPayFrequencyNotFound = "No pay frequency were found."
)

// GetPayFrequency gets the information of table pay frequency from database
// GetPayFrequency returns all rows of Pay Frequency if the id parameter has no value, otherwise single row will be returned
func (co *Controllers) GetPayFrequency(c echo.Context) error {
	// get query pay frequecy id
	queryPayFrequecyID := c.QueryParam("id")

	// If the query pay frequency id is empty, then it displays all information of table pay frequency from database
	if len(queryPayFrequecyID) == 0 {
		payFrequencies := PayFrequencies{}

		err := db.Select().
			From("tblpayfrequency").
			All(&payFrequencies)

		if err != nil {
			message := ErrPayFrequencyNotFound
			return SendErrorResponse(c, http.StatusNotFound, message)
		}

		return SendResponse(c, http.StatusOK, payFrequencies)
	}

	// payFrequencyID is converted to integer.
	payFrequencyID, err := strconv.Atoi(queryPayFrequecyID)
	if err != nil {
		message := ErrPayFrequencyNotFound
		return SendErrorResponse(c, http.StatusNotFound, message)
	}

	// Display all information of table payfrequency from database if it is exist
	payFrequency := PayFrequency{}
	err = db.Select().
		From("tblpayfrequency").
		Where(dbx.HashExp{"id": payFrequencyID}).
		One(&payFrequency)
	if err != nil {
		message := ErrPayFrequencyNotFound
		return SendErrorResponse(c, http.StatusNotFound, message)
	}

	return SendResponse(c, http.StatusOK, payFrequency)

}
