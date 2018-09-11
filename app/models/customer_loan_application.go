package models

import null "gopkg.in/guregu/null.v3"

// CustomerLoanApplication contains information about a loan application
type CustomerLoanApplication struct {
	ID                  int
	CustomerID          null.Int `db:"fk_customer_id"`
	LoanAmount          null.Float
	InterestAmount      null.Float
	FeeAmount           null.Float
	TotalAmount         null.Float
	ReferenceCode       null.String
	DueDate             null.String
	LoanDate            null.String
	CreatedBy           null.String
	CreatedDate         null.String
	Status              null.String
	ProcessedBy         null.String
	ProcessedDate       null.String
	Source              null.String
	EpointTransactionID null.String
}

// TableName gets the name of the database table
func (cla CustomerLoanApplication) TableName() string {
	return "tblCustomerLoanApplication"
}
