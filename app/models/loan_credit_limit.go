package models

import (
	"time"

	null "gopkg.in/guregu/null.v3"
)

// LoanCreditLimit contains information about a credit limit
type LoanCreditLimit struct {
	ID           int
	Tier         null.Int
	Code         null.String
	Description  null.String
	Amount       null.Float
	NumberOfDays null.Int `db:"no_of_days"`
	Active       null.String
	CreatedDate  time.Time
}

// TableName gets the name of the database table
func (lcl LoanCreditLimit) TableName() string {
	return "tblLoanCreditLimit"
}
