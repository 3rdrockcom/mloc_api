package controllers

import (
	"net/http"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/labstack/echo"
)

// CustomerInfos array contains customerInfo struct
type CustomerInfos []CustomerInfo // this is for GetCustomerLoanList function

// CustomerInfo contains information about fkCustomerId
type CustomerInfo struct {
	CustUniqueID     *string  `db:"cust_unique_id" json:"cust_unique_id"`
	CustomerID       int      `db:"customer_id" json:"customer_id"`
	FirstName        *string  `db:"first_name" json:"first_name"`
	MiddleName       *string  `db:"middle_name" json:"middle_name"`
	LastName         *string  `db:"last_name" json:"last_name"`
	Suffix           *string  `db:"suffix" json:"suffix"`
	BirthDate        *string  `db:"birth_date" json:"birth_date"` //may need to change date type
	Address1         *string  `db:"address1" json:"address1"`
	Address2         *string  `db:"address2" json:"address2"`
	CountryID        *int     `db:"country_id" json:"country_id"`
	CountryDesc      *string  `db:"country_desc" json:"country_desc"`
	StateID          *int     `db:"state_id" json:"state_id"`
	StateDesc        *string  `db:"state_desc" json:"state_desc"`
	CityID           *int     `db:"city_id" json:"city_id"`
	CityDesc         *string  `db:"city_desc" json:"city_desc"`
	ZipCode          *string  `db:"zipcode" json:"zipcode"`
	HomeNumber       *string  `db:"home_number" json:"home_number"`
	MobileNumber     *string  `db:"mobile_number" json:"mobile_number"`
	Email            *string  `db:"email" json:"email"`
	CompanyName      *string  `db:"company_name" json:"company_name"`
	PhoneNumber      *string  `db:"phone_number" json:"phone_number"`
	NetPayPerCheck   *float64 `db:"net_pay_percheck" json:"net_pay_percheck"`
	IncomeSourceID   *int     `db:"income_source_id" json:"income_source_id"`
	MlocAccess       *int     `db:"mloc_access" json:"mloc_access"`
	Registration     *int     `db:"registration" json:"registration"`
	TermAndCondition *int     `db:"term_and_condition" json:"term_and_condition"`
	IncomeSourceDesc *string  `db:"income_source_desc" json:"income_source_desc"`
	PayFrequencyID   *int     `db:"pay_frequency_id" json:"pay_frequency_id"`
	PayFrequencyDesc *string  `db:"pay_frequency_desc" json:"pay_frequency_desc"`
	NextPayDate      *string  `db:"next_paydate" json:"next_paydate"` //may need to chage date type
	Key              *string  `db:"key" json:"key"`
	CreditLimit      *float64 `db:"credit_limit" json:"credit_limit"`
	AvailableCredit  *float64 `db:"available_credit" json:"available_credit"`
	IsSuspended      *string  `db:"is_suspended" json:"is_suspended"` // may need to change enum type
}

var (
	// ErrCustomerNotFound is error given when country is not found
	ErrCustomerNotFound = "No customer was found."
)

// GetCustomerInfo displays all transaction history of customer
func (co *Controllers) GetCustomerInfo(c echo.Context) error {

	// Get customer ID
	fkCustomerID := c.Get("customerID").(int)

	customerInfo := CustomerInfo{}

	err := db.Select().
		From("view_customer_info").
		Where(dbx.HashExp{"customer_id": fkCustomerID}).
		OrderBy("customer_id").
		One(&customerInfo)

	if err != nil {
		message := ErrCustomerNotFound
		return SendErrorResponse(c, http.StatusNotFound, message)

	}
	if customerInfo.IsSuspended != nil {
		if *customerInfo.IsSuspended == "NO" {
			*customerInfo.IsSuspended = "0"
		}

		if *customerInfo.IsSuspended == "YES" {
			*customerInfo.IsSuspended = "1"
		}
	}

	return SendResponse(c, http.StatusOK, customerInfo)
}
