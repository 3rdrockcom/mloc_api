package models

import null "gopkg.in/guregu/null.v3"

// CustomerPayments is an array of customer payments
type CustomerPayments []CustomerPayment

// CustomerPayment contains information about a customer's payments
type CustomerPayment struct {
	ID                  int
	CustomerID          null.Int `db:"fk_customer_id"`
	ReferenceCode       null.String
	EpointTransactionID null.String
	PaymentAmount       null.Float
	DatePaid            null.String
	PaidBy              null.String
}

// TableName gets the name of the database table
func (cp CustomerPayment) TableName() string {
	return "tblCustomerPayment"
}
