package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// Customers is an array of Customer entries
type Customers []Customer

// Customer contains information about a customer
type Customer struct {
	ID                    int    `json:"id"`
	FirstName             string `json:"first_name" binding:"required"`
	LastName              string `json:"last_name" binding:"required"`
	Email                 string `json:"email" binding:"required"`
	ProgramID             int    `json:"-"`
	ProgramCustomerID     int    `json:"-"`
	ProgramCustomerMobile string `json:"-"`
	CustomerUniqueID      string `json:"-" db:"cust_unique_id"`
}

// TableName gets the name of the database table
func (c Customer) TableName() string {
	return "tblcustomerbasicinfo"
}

// Validate checks if the values in the struct are valid
func (c Customer) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.FirstName, validation.Required),
		validation.Field(&c.LastName, validation.Required),
		validation.Field(&c.Email, validation.Required, is.Email),
	)
}
