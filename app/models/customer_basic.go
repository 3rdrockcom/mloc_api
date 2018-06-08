package models

import null "gopkg.in/guregu/null.v3"

// CustomerBasic contains basic information of a customer
type CustomerBasic struct {
	ID           int
	FirstName    null.String `db:"first_name"`
	MiddleName   null.String `db:"middle_name"`
	LastName     null.String `db:"last_name"`
	Suffix       null.String `db:"suffix"`
	Birthday     null.String `db:"birth_date"` // may need to change date type
	Address1     null.String `db:"address1"`
	Address2     null.String `db:"address2"`
	Country      null.Int    `db:"country"`
	State        null.Int    `db:"state"`
	City         null.Int    `db:"city"`
	ZipCode      null.String `db:"zipcode"`
	HomeNumber   null.String `db:"home_number"`
	MobileNumber null.String `db:"mobile_number"`
	Email        null.String `db:"email"`
}

// TableName gets customerbasicinfo table from database
func (a CustomerBasic) TableName() string {
	return "tblCustomerBasicInfo"
}
