package models

import (
	"gopkg.in/guregu/null.v3"
)

// CustomerAdditional contains information of customer id
type CustomerAdditional struct {
	ID               int
	CustomerID       null.Int    `db:"fk_customer_id"`
	CompanyName      null.String `db:"company_name"`
	PhoneNumber      null.String `db:"phone_number"`
	NetPayPerCheck   null.Float  `db:"net_pay_percheck"`
	IncomeSource     null.Int    `db:"income_source"`
	PayFrequency     null.Int    `db:"pay_frequency"`
	NextPayDate      null.String `db:"next_paydate"`      // may need to change date type
	FollowingPayDate null.String `db:"following_paydate"` // may need to change date type
}

// TableName gets tblcustomerotherinfo table from database
func (ca CustomerAdditional) TableName() string {
	return "tblCustomerOtherInfo"
}
