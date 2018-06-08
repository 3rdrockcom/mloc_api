package controllers

import (
	"net/http"
	"strconv"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/labstack/echo"
)

// CustAdditional contains information of customer id
type CustAdditional struct {
	ID               int
	FkCustomerID     *int     `db:"fk_customer_id"`
	CompanyName      *string  `db:"company_name"`
	PhoneNumber      *string  `db:"phone_number"`
	NetPayPercheck   *float64 `db:"net_pay_percheck"`
	IncomeSource     *int     `db:"income_source"`
	Payfrequency     *int     `db:"pay_frequency"`
	NextPayDate      *string  `db:"next_paydate"`      // may need to change date type
	FollowingPayDate *string  `db:"following_paydate"` // may need to change date type
	//CustUniqueId *string `db:"cust_unique_id`
}

// TableName gets tblcustomerotherinfo table from database
func (c CustAdditional) TableName() string {
	return "tblcustomerotherinfo"
}

var (
	// ErrIncompleteCustInfo is given if the customer information is incomplete
	ErrIncompleteCustInfo = "Provide complete customer information to create."

	// MsgCustomerInfoUpdateSuccess is given if the customer information update successfully
	MsgCustomerInfoUpdateSuccess = "Customer information has been updated successfully."
)

// PostCustomerAdditional update customer information in tblcustomerotherinfo from database
func (co *Controllers) PostCustomerAdditional(c echo.Context) error {

	// Get customerID
	customerID := c.Get("customerID").(int)
	customerOtherinfo := &CustAdditional{}

	err := db.Select().
		From("tblcustomerotherinfo").
		Where(dbx.HashExp{"Fk_customer_id": customerID}).
		One(customerOtherinfo)

	if err != nil {
		message := ErrIncompleteCustInfo
		return SendErrorResponse(c, http.StatusBadRequest, message)
	}

	formKeys, err := c.FormParams()

	for formKey := range formKeys {
		value := c.FormValue(formKey) // read each value from postform

		switch formKey {
		case "R1":
			customerOtherinfo.CompanyName = &value

		case "R2":
			customerOtherinfo.PhoneNumber = &value

		case "R3":
			fvalue, err := strconv.ParseFloat(value, 64)
			if err != nil {
				message := ErrIncompleteCustInfo
				return SendErrorResponse(c, http.StatusBadRequest, message)
			}
			customerOtherinfo.NetPayPercheck = &fvalue

		case "R4":
			var tempVal int
			tempVal, err = strconv.Atoi(value)
			if err != nil {
				message := ErrIncompleteCustInfo
				return SendErrorResponse(c, http.StatusBadRequest, message)
			}
			customerOtherinfo.IncomeSource = &tempVal

		case "R5":
			var tempVal int
			tempVal, err = strconv.Atoi(value)
			if err != nil {
				message := ErrIncompleteCustInfo
				return SendErrorResponse(c, http.StatusBadRequest, message)
			}
			customerOtherinfo.Payfrequency = &tempVal

		case "R6":
			customerOtherinfo.NextPayDate = &value

		case "R7":
			customerOtherinfo.FollowingPayDate = &value

		case "R8":
			// Not necessary to do
		}

	}

	err = db.Model(customerOtherinfo).Update()
	if err != nil {
		if err != nil {
			message := ErrIncompleteCustInfo
			return SendErrorResponse(c, http.StatusBadRequest, message)
		}
	} // end of valid postform
	message := MsgCustomerInfoUpdateSuccess
	return SendOKResponse(c, message)

}
