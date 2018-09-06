package controllers

import (
	"net/http"

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
