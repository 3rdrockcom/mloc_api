package controllers

import (
	"net/http"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/labstack/echo"
)

// LoanHistories is array of loanHistory type
type LoanHistories []LoanHistory

// LoanHistory contains information about fkCustomerId
type LoanHistory struct {
	FkCustomerID *int     `db:"fk_customer_id" json:"fk_customer_id"`
	Amount       *float64 `db:"amount" json:"amount"`
	Ttype        string   `db:"t_type" json:"t_type"`
	TDate        *string  `db:"t_date" json:"t_date"` //may need to change datetime type
}

var (
	// ErrTransactionNotFound is error given when country is not found
	ErrTransactionNotFound = "No transaction(s) were found."
)

// GetTransactionHistory displays all transaction history of customer base on transaction type
func (co *Controllers) GetTransactionHistory(c echo.Context) error {

	fkCustomerID := c.Get("customerID").(int)

	// Get transaction type
	queryLoanOrPayment := c.QueryParam("R2")
	loanHistories := LoanHistories{}

	// If transaction type is LOAN Or PAYMENT, it displays all loan or payment for customer
	if len(queryLoanOrPayment) > 0 {
		if queryLoanOrPayment == "LOAN" || queryLoanOrPayment == "PAYMENT" {
			err := db.Select().
				From("view_transaction_history").
				Where(dbx.HashExp{"t_type": queryLoanOrPayment, "fk_customer_id": fkCustomerID}).
				OrderBy("fk_customer_id").
				All(&loanHistories)

			if err != nil {
				message := ErrTransactionNotFound

				return SendErrorResponse(c, http.StatusNotFound, message)
			}

			if len(loanHistories) == 0 {
				message := ErrTransactionNotFound
				return SendErrorResponse(c, http.StatusNotFound, message)
			}

			return SendResponse(c, http.StatusOK, loanHistories)
		}
	}

	// It displays all loan and payment for customer if the transaction type is empty
	err := db.Select().
		From("view_transaction_history").
		Where(dbx.HashExp{"fk_customer_id": fkCustomerID}).
		OrderBy("fk_customer_id").
		All(&loanHistories)

	if err != nil {
		message := ErrTransactionNotFound
		return SendErrorResponse(c, http.StatusNotFound, message)
	}

	if len(loanHistories) == 0 {
		message := ErrTransactionNotFound
		return SendErrorResponse(c, http.StatusNotFound, message)
	}

	return SendResponse(c, http.StatusOK, loanHistories)

}
