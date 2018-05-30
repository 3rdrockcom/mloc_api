package controllers

import (
	"net/http"
	"strconv"

	dbx "github.com/go-ozzo/ozzo-dbx"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/labstack/echo"
)

// PostCustomerBasic contains information of table customer basic
type PostCustomerBasic struct {
	ID           int
	FirstName    *string `db:"first_name"`
	MiddleName   *string `db:"middle_name"`
	LastName     *string `db:"last_name"`
	Suffix       *string `db:"suffix"`
	Birthday     *string `db:"birth_date"` // may need to change date type
	Address1     *string `db:"address1"`
	Address2     *string `db:"address2"`
	Country      *int    `db:"country"`
	State        *int    `db:"state"`
	City         *int    `db:"city"`
	ZipCode      *string `db:"zipcode"`
	HomeNumber   *string `db:"home_number"`
	MobileNumber *string `db:"mobile_number"`
	Email        *string `db:"email"`
}

// TableName gets customerbasicinfo table from database
func (a PostCustomerBasic) TableName() string {
	return "tblcustomerbasicinfo"
}

// Validate checks postform required is validation
func (a PostCustomerBasic) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.FirstName, validation.Required),
		validation.Field(&a.LastName, validation.Required),
		validation.Field(&a.MobileNumber, validation.Required),
		validation.Field(&a.Email, validation.Required, is.Email),
	)
}

var (
	// ErrProblemOccured is given if it can't get data from database or can't covert input to string
	ErrProblemOccured = "Some problems occurred, please try agin."

	// ErrCustomerInfoIncomplete is given if customerinformation is not complete
	ErrCustomerInfoIncomplete = "Provide complete customer information to create."

	// MsgCustomerInfoSuccess is given if customer information is update successfully
	MsgCustomerInfoSuccess = "Customer information has been updated successfully."
)

// PostCustomerBasic allows to update table customer basic information from database if it is valid
func (co *Controllers) PostCustomerBasic(c echo.Context) error {
	// Get customer ID
	customerID := c.Get("customerID").(int)

	customerBasic := &PostCustomerBasic{}
	customerBasicRequired := PostCustomerBasic{}

	// Get table customer basic information from database
	err := db.Select().
		From("tblcustomerbasicinfo").
		Where(dbx.HashExp{"id": customerID}).
		One(customerBasic)

	// It displays error if the customer id is not exists
	if err != nil {
		message := ErrProblemOccured
		return SendErrorResponse(c, http.StatusBadRequest, message)
	}

	formKeys, err := c.FormParams()

	for formKey := range formKeys {
		value := c.FormValue(formKey) // read each value from postform

		switch formKey {
		case "R1":
			customerBasic.FirstName = &value
			customerBasicRequired.FirstName = &value

		case "R2":
			customerBasic.MiddleName = &value

		case "R3":
			customerBasic.LastName = &value
			customerBasicRequired.LastName = &value

		case "R4":
			customerBasic.Suffix = &value
		case "R5":
			customerBasic.Birthday = &value

		case "R6":
			customerBasic.Address1 = &value

		case "R7":
			customerBasic.Address2 = &value

		case "R8":
			if len(value) == 0 {
				customerBasic.Country = nil
				continue
			}
			var tempVal int
			tempVal, err = strconv.Atoi(value)
			if err != nil {

				message := ErrProblemOccured
				return SendErrorResponse(c, http.StatusBadRequest, message)
			}
			customerBasic.Country = &tempVal

		case "R9":
			if len(value) == 0 {
				customerBasic.State = nil
				continue
			}
			var tempVal int
			tempVal, err = strconv.Atoi(value)
			if err != nil {
				message := ErrProblemOccured
				return SendErrorResponse(c, http.StatusBadRequest, message)
			}

			customerBasic.State = &tempVal

		case "R10":
			if len(value) == 0 {
				customerBasic.City = nil
				continue
			}

			var tempVal int
			tempVal, err = strconv.Atoi(value)
			if err != nil {
				message := ErrProblemOccured
				return SendErrorResponse(c, http.StatusBadRequest, message)
			}

			customerBasic.City = &tempVal

		case "R11":
			customerBasic.ZipCode = &value

		case "R12":
			customerBasic.HomeNumber = &value

		case "R13":
			customerBasic.MobileNumber = &value
			customerBasicRequired.MobileNumber = &value

		case "R14":
			customerBasic.Email = &value
			customerBasicRequired.Email = &value
		case "R15":
			// Not necessary to do
		case "R16":
			//Not neccessary to do
		}
	}

	// Check customer required is validate
	err = customerBasicRequired.Validate()
	if err != nil {
		message := ErrCustomerInfoIncomplete
		return SendErrorResponse(c, http.StatusBadRequest, message)
	}
	// Check require postform is valid or not
	err = customerBasic.Validate()

	// Require postform is not validation, displays error
	if err != nil {
		message := ErrCustomerInfoIncomplete
		return SendErrorResponse(c, http.StatusNotFound, message)
	}

	// Displays update customer error occur
	err = db.Model(customerBasic).Update()
	if err != nil {
		message := ErrProblemOccured
		return SendErrorResponse(c, http.StatusBadRequest, message)
	}

	// Success to update customer basic
	message := MsgCustomerInfoSuccess
	return SendResponse(c, http.StatusOK, message)

}
