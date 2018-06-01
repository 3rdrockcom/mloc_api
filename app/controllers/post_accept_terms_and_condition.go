package controllers

import (
	"fmt"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/labstack/echo"
)

// CustomerAgreement contains information of customer id
type CustomerAgreement struct {
	ID               int
	TermAndCondition *int `db:"term_and_condition"`
}

// TableName gets customerbasicinfo table from database
func (a CustomerAgreement) TableName() string {
	return "tblcustomeragreement"
}

var (
	// MsgCustomerAgreef is ginven if the customer agree in terms and condition
	MsgCustomerAgreef = "Customer agreed in terms and condition."
)

// PostAcceptTermsAndCondition accepts the term and condition
// the value of accept will store in tblcustomeragreement
func (co *Controllers) PostAcceptTermsAndCondition(c echo.Context) error {

	customerID := c.Get("customerID").(int)
	customerAgreement := &CustomerAgreement{}

	err := db.Select().
		From("tblcustomeragreement").
		Where(dbx.HashExp{"fk_customer_id": customerID}).
		One(customerAgreement)
	if err != nil {
		return nil
	}

	*customerAgreement.TermAndCondition = 1
	fmt.Println(*customerAgreement)

	err = db.Model(customerAgreement).Update()
	if err != nil {
		return nil
	}

	message := MsgCustomerInfoUpdateSuccess
	return SendOKResponse(c, message)

}
