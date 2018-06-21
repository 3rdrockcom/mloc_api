package models

import (
	"gopkg.in/guregu/null.v3"
)

// LoanTotal contains information about a loan totals
type LoanTotal struct {
	ID                   int        `db:"id" json:"id"`
	CustomerID           null.Int   `db:"fk_customer_id" json:"fk_customer_id"`
	TotalPrincipalAmount null.Float `db:"total_principal_amount" json:"-"`
	TotalFeeAmount       null.Float `db:"total_fee_amount" json:"-"`
	TotalAmount          null.Float `db:"total_amount" json:"-"`

	TotalPrincipalAmountDisplay null.String `db:"-" json:"total_principal_amount"`
	TotalFeeAmountDisplay       null.String `db:"-" json:"total_fee_amount"`
	TotalAmountDisplay          null.String `db:"-" json:"total_amount"`
}

// TableName gets the name of the database table
func (lt LoanTotal) TableName() string {
	return "tblCustomerLoanTotal"
}
