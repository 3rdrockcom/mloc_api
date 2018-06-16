package models

import (
	null "gopkg.in/guregu/null.v3"
	"gopkg.in/guregu/null.v3/zero"
)

// CustomerInfo contains detailed information about a customer
type CustomerInfo struct {
	ID                    int         `db:"customer_id" json:"customer_id"`
	CustUniqueID          null.String `db:"cust_unique_id" json:"cust_unique_id"`
	FirstName             null.String `db:"first_name" json:"first_name"`
	MiddleName            null.String `db:"middle_name" json:"middle_name"`
	LastName              null.String `db:"last_name" json:"last_name"`
	Suffix                null.String `db:"suffix" json:"suffix"`
	BirthDate             null.String `db:"birth_date" json:"birth_date"` //may need to change date type
	Address1              null.String `db:"address1" json:"address1"`
	Address2              null.String `db:"address2" json:"address2"`
	CountryID             null.Int    `db:"country_id" json:"country_id"`
	CountryDesc           null.String `db:"country_desc" json:"country_desc"`
	StateID               null.Int    `db:"state_id" json:"state_id"`
	StateDesc             null.String `db:"state_desc" json:"state_desc"`
	CityID                null.Int    `db:"city_id" json:"city_id"`
	CityDesc              null.String `db:"city_desc" json:"city_desc"`
	ZipCode               null.String `db:"zipcode" json:"zipcode"`
	HomeNumber            null.String `db:"home_number" json:"home_number"`
	MobileNumber          null.String `db:"mobile_number" json:"mobile_number"`
	Email                 null.String `db:"email" json:"email"`
	CompanyName           null.String `db:"company_name" json:"company_name"`
	PhoneNumber           null.String `db:"phone_number" json:"phone_number"`
	NetPayPerCheck        null.Float  `db:"net_pay_percheck" json:"net_pay_percheck"`
	IncomeSourceID        null.Int    `db:"income_source_id" json:"income_source_id"`
	MLOCAccess            zero.Int    `db:"mloc_access" json:"mloc_access"`
	Registration          zero.Int    `db:"registration" json:"registration"`
	TermsAndConditions    zero.Int    `db:"term_and_condition" json:"term_and_condition"`
	IncomeSourceDesc      null.String `db:"income_source_desc" json:"income_source_desc"`
	PayFrequencyID        null.Int    `db:"pay_frequency_id" json:"pay_frequency_id"`
	PayFrequencyDesc      null.String `db:"pay_frequency_desc" json:"pay_frequency_desc"`
	NextPayDate           null.String `db:"next_paydate" json:"next_paydate"` //may need to chage date type
	Key                   null.String `db:"key" json:"key"`
	CreditLimit           null.Float  `db:"credit_limit" json:"credit_limit"`
	AvailableCredit       null.Float  `db:"available_credit" json:"available_credit"`
	IsSuspended           string      `db:"is_suspended" json:"is_suspended"`
	CreditLineID          null.Int    `json:"-"`
	ProgramCustomerID     null.Int    `json:"-"`
	ProgramCustomerMobile null.String `json:"-"`
}

// TableName gets the name of the database table
func (c CustomerInfo) TableName() string {
	return "view_customer_info"
}
