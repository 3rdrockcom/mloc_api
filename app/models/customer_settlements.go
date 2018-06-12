package models

import (
	"gopkg.in/guregu/null.v3"
)

// CustomerSettlements is an array of customer settlements
type CustomerSettlements []CustomerSettlement

// CustomerSettlement contains information about a customer's settlement
type CustomerSettlement struct {
	ID                int
	CustomerID        null.Int `db:"fk_customer_id"`
	CustomerLoanID    null.Int
	CustomerPaymentID null.Int
	SettlementAmount  null.Float
	PrincipalAmount   null.Float
	FeeAmount         null.Float
	CreatedDate       null.String
}

// TableName gets the name of the database table
func (cs CustomerSettlement) TableName() string {
	return "tblCustomerSettlement"
}
