package controllers

import (
	"fmt"
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
	ProgramID    *int    `db:"program_id"`
	CustUniqueID *string `db:"cust_unique_id"`
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
	// ProblemOccurredMessage displays error message occur when update table customer basic information fails
	ProblemOccurredMessage = "Some problems occurred, please try agin."
)

// PostCustomerBasic allows to update table customer basic information from database if it is valid
func (co *Controllers) PostCustomerBasic(c echo.Context) error {
	// Get customer ID
	customerID := c.Get("customerID").(int)

	customerBasic := &PostCustomerBasic{}

	// get table customer basic information from database
	err := db.Select().
		From("tblcustomerbasicinfo").
		Where(dbx.HashExp{"id": customerID}).
		One(customerBasic)

		// it displays error if the customer id is not exists
	if err != nil {
		message := ProblemOccurredMessage
		fmt.Println("not get it from db")
		return SendErrorResponse(c, http.StatusNotFound, message)
	}

	// array of postform
	formKeys := []string{"R1", "R2", "R3", "R4", "R5", "R6", "R7", "R8", "R9", "R10", "R11", "R12", "R13", "R14", "R15", "R16"}

	for index := range formKeys {
		formKey := formKeys[index]
		value := c.FormValue(formKey) // read each value from postform

		switch formKey {
		case "R1":
			customerBasic.FirstName = &value
		case "R2":
			customerBasic.MiddleName = &value
		case "R3":
			customerBasic.LastName = &value
		case "R4":
			customerBasic.Suffix = &value
		case "R5":
			customerBasic.Birthday = &value
		case "R6":
			customerBasic.Address1 = &value
		case "R7":
			customerBasic.Address2 = &value

		case "R8":
			var tempVal int
			tempVal, err = strconv.Atoi(value)
			if err != nil {

				message := ProblemOccurredMessage
				fmt.Println("checking end R8")
				return SendErrorResponse(c, http.StatusNotFound, message)
			}
			customerBasic.Country = &tempVal
		case "R9":

			var tempVal int
			tempVal, err = strconv.Atoi(value)
			if err != nil {
				message := ProblemOccurredMessage
				return SendErrorResponse(c, http.StatusNotFound, message)
			}

			customerBasic.State = &tempVal
		case "R10":

			var tempVal int
			tempVal, err = strconv.Atoi(value)
			if err != nil {
				message := ProblemOccurredMessage
				return SendErrorResponse(c, http.StatusNotFound, message)
			}

			customerBasic.City = &tempVal
		case "R11":
			fmt.Println("checking begin R11")
			customerBasic.ZipCode = &value
		case "R12":
			fmt.Println("checking begin R12")
			customerBasic.HomeNumber = &value
		case "R13":
			fmt.Println("checking begin R13")
			customerBasic.MobileNumber = &value
		case "R14":
			fmt.Println("checking begin R14")
			customerBasic.Email = &value
		case "R15": // not require
			/*
				var tempVal int
				tempVal, err = strconv.Atoi(value)
				if err != nil {
					message := ProblemOccurredMessage
					fmt.Println("checking error in R15")
					return SendErrorResponse(c, http.StatusNotFound, message)
				}
				customerBasic.ProgramID = &tempVal
			*/
		case "R16":

			// not required
			// customerBasic.CustUniqueID = &value

			//	}

		}
	}
	//check require postform is valid or not
	err = customerBasic.Validate()

	// require postform is not validation, displays error
	if err != nil {
		ProblemOccurredMessage := "Provide complete customer information to create."
		fmt.Println("not get it to check validate")
		return SendErrorResponse(c, http.StatusNotFound, ProblemOccurredMessage)
	}
	fmt.Println("checking after validate")
	// displays update customer error occur
	err = db.Model(customerBasic).Update()
	if err != nil {
		message := ProblemOccurredMessage
		return SendErrorResponse(c, http.StatusNotFound, message)
	}

	//success to update customer basic
	successUpdataMessage := "Customer information has been updated successfully."
	return SendResponse(c, http.StatusOK, successUpdataMessage)

}
