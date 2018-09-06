package controllers

import (
	"net/http"

	Customer "github.com/epointpayment/mloc_api_go/app/services/customer"
	"github.com/epointpayment/mloc_api_go/app/services/payments"
	"github.com/epointpayment/mloc_api_go/app/services/payments/institution"

	Payments "github.com/epointpayment/mloc_api_go/app/services/payments"

	"github.com/labstack/echo"
)

//
func (co Controllers) GetDisbursementMethod(c echo.Context) (err error) {
	// Initialize lookup service
	ls := Payments.New()

	// Get list of disbursement methods
	methods, err := ls.GetDisbursementMethods()
	if err != nil {
		return err
	}

	// Send response
	return SendResponse(c, http.StatusOK, methods)
}

//
func (co Controllers) GetCollectionMethod(c echo.Context) (err error) {
	// Initialize lookup service
	ls := Payments.New()

	// Get list of collection methods
	methods, err := ls.GetCollectionMethods()
	if err != nil {
		return err
	}

	// Send response
	return SendResponse(c, http.StatusOK, methods)
}

//
func (co Controllers) GetBankInstitutions(c echo.Context) (err error) {
	// Get customer ID
	customerID := c.Get("customerID").(int)

	// Initialize customer service
	sc, err := Customer.New(customerID)
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// Get detailed customer information
	customerInfo, err := sc.Info().GetDetails()
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// Initialize lookup service
	ls := Payments.New()

	// Get list of collection methods
	institutions, err := ls.GetInstitutions(institution.Request{
		Method:   payments.MethodSTP,
		Customer: *customerInfo,
	})
	if err != nil {
		return err
	}

	// Send response
	return SendResponse(c, http.StatusOK, institutions.Institutions)
}
