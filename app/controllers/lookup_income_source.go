package controllers

import (
	"net/http"
	"strconv"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/labstack/echo"
)

// IncomeSources is an array of IncomeSouce entries
type IncomeSources []IncomeSource

// IncomeSource contains information about a IncomeSourceID
type IncomeSource struct {
	ID          int     `db:"id" json:"id"`
	Description *string `db:"description" json:"description"`
}

// TableName gets the name of the database table
func (c IncomeSource) TableName() string {
	return "tblincomesource"
}

var (
	// ErrIncomeSourceNotFound is error given when income source is not found
	ErrIncomeSourceNotFound = "No income source was found."
)

// GetIncomeSource gets the information of table income source from database
// returns all rows of income source if the id parameter has no value, otherwise single row will be returned
func (co *Controllers) GetIncomeSource(c echo.Context) error {
	// Get query income source id
	queryIncomeSourceID := c.QueryParam("id")

	// If the query income source id is empty, then it displays all information of table income source from database
	if len(queryIncomeSourceID) == 0 {
		incomeSources := IncomeSources{}

		err := db.Select().
			From("tblincomesource").
			All(&incomeSources)
		if err != nil {
			message := ErrIncomeSourceNotFound
			return SendErrorResponse(c, http.StatusNotFound, message)
		}

		return SendResponse(c, http.StatusOK, incomeSources)
	}

	// queryIncomeSourceId is converted to integer.
	incomeSourceID, err := strconv.Atoi(queryIncomeSourceID)
	if err != nil {
		message := ErrIncomeSourceNotFound
		return SendErrorResponse(c, http.StatusNotFound, message)
	}

	// Display all information of table income source from database if it is exist
	incomeSource := IncomeSource{}
	err = db.Select().
		From("tblincomesource").
		Where(dbx.HashExp{"id": incomeSourceID}).
		One(&incomeSource)
	if err != nil {
		message := ErrIncomeSourceNotFound
		return SendErrorResponse(c, http.StatusNotFound, message)
	}

	return SendResponse(c, http.StatusOK, incomeSource)

}
