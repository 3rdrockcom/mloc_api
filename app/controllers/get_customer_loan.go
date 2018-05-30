package controllers

import (
	"net/http"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/labstack/echo"
)

// Loan contains information about a Loan ID
type Loan struct {
	ID                   int      `db:"id" json:"id"`
	FKCustomerID         *int     `db:"fk_customer_id" json:"fk_customer_id"`
	TotalPrincipalAmount *float64 `db:"total_principal_amount" json:"total_principal_amount"`
	TotalFeeAmount       *float64 `db:"total_fee_amount" json:"total_fee_amount"`
	TotalAmount          *float64 `db:"total_amount" json:"total_amount"`
}

var (
	// ErrLoanNotFound is error given when customer loan is not found
	ErrLoanNotFound = "No customer loan was found."
)

// GetCustomerLoan displays information of customer loan if the customer exists
func (co *Controllers) GetCustomerLoan(c echo.Context) error {
	// Get fk customer id
	fkCustomerID := c.Get("customerID").(int)

	loan := Loan{}

	//Displays the customer information if fk_customer_id is exists
	err := db.Select().
		From("tblCustomerLoanTotal").
		Where(dbx.HashExp{"fk_customer_id": fkCustomerID}).
		OrderBy("fk_customer_id").
		One(&loan)

	if err != nil {
		message := ErrLoanNotFound
		return SendErrorResponse(c, http.StatusNotFound, message)
	}

	return SendResponse(c, http.StatusOK, loan)

}
