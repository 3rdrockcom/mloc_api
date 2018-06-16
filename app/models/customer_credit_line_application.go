package models

import (
	"time"

	null "gopkg.in/guregu/null.v3"
)

// CustomerCreditLineApplication contains information about a customer's credit line
type CustomerCreditLineApplication struct {
	ID               int
	CustomerID       null.Int `db:"fk_customer_id"`
	CreditLineID     null.Int
	CreditLineAmount null.Float
	ReferenceCode    null.String
	Status           null.String
	ProcessedBy      null.String
	ProcessedDate    null.String // may need to date type
	CreatedDate      time.Time
}

// TableName gets the name of the database table
func (ccla CustomerCreditLineApplication) TableName() string {
	return "tblCustomerCreditLineApplication"
}
