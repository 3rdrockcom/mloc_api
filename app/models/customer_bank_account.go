package models

import "time"

type CustomerBankAccount struct {
	ID            int       `json:"id"`
	CustomerID    int       `json:"-"`
	Alias         string    `json:"alias"`
	BankCode      string    `json:"bank_code"`
	AccountNumber string    `json:"account_number"`
	DateCreated   time.Time `json:"-"`
	DateUpdated   time.Time `json:"-"`

	DateCreatedDisplay string `db:"-" json:"date_created"`
	DateUpdatedDisplay string `db:"-" json:"date_updated"`
}

// TableName gets the name of the database table
func (m CustomerBankAccount) TableName() string {
	return "tblCustomerBankAccount"
}
