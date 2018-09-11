package models

import (
	"time"

	"gopkg.in/guregu/null.v3"
)

// CustomerLoans is an array of customer loans
type CustomerLoans []CustomerLoan

// CustomerLoan contains information about a customer loan
type CustomerLoan struct {
	ID                  int
	CustomerID          null.Int `db:"fk_customer_id"`
	LoanApplicationID   null.Int
	Source              null.String
	EpointTransactionID null.String
	LoanIntervalID      null.Int
	LoanTermID          null.Int
	LoanAmount          null.Float
	InterestAmount      null.Float
	FeeAmount           null.Float
	TotalAmount         null.Float
	TotalPaidPrincipal  null.Float
	TotalPaidFee        null.Float
	TotalPaidAmount     null.Float
	ReferenceCode       null.String
	IsPaid              null.Int
	DueDate             null.String
	LoanDate            null.String
	ApprovedBy          null.String
	ApprovedDate        null.String
	CreatedDate         time.Time
}

// TableName gets the name of the database table
func (cl CustomerLoan) TableName() string {
	return "tblCustomerLoan"
}
