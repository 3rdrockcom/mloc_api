package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	null "gopkg.in/guregu/null.v3"
)

// Customers is an array of Customer entries
type Customers []Customer

// Customer contains information about a customer
type Customer struct {
	ID                    int         `json:"id"`
	FirstName             null.String `json:"first_name" binding:"required"`
	LastName              null.String `json:"last_name" binding:"required"`
	Email                 null.String `json:"email" binding:"required"`
	ProgramID             null.Int    `json:"-"`
	ProgramCustomerID     null.Int    `json:"-"`
	ProgramCustomerMobile null.String `json:"-"`
	CustomerUniqueID      null.String `json:"-" db:"cust_unique_id"`
}

// TableName gets the name of the database table
func (c Customer) TableName() string {
	return "tblCustomerBasicInfo"
}

// Validate checks if the values in the struct are valid
func (c Customer) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.FirstName, validation.Required),
		validation.Field(&c.LastName, validation.Required),
		validation.Field(&c.Email, validation.Required, is.Email),
	)
}
