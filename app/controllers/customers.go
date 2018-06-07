package controllers

import (
	"net/http"

	Customer "github.com/epointpayment/mloc_api_go/app/services/customer"

	"github.com/labstack/echo"
)

// GetCustomer displays detailed customer information
func (co *Controllers) GetCustomer(c echo.Context) error {
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

	switch customerInfo.IsSuspended {
	case "YES":
		customerInfo.IsSuspended = "1"
	case "NO":
		fallthrough
	default:
		customerInfo.IsSuspended = "0"
	}

	return SendResponse(c, http.StatusOK, customerInfo)
}

// PostAddCustomer updates customer information
func (co Controllers) PostAddCustomer(c echo.Context) error {
	// Get customer ID
	customerID := c.Get("customerID").(int)

	// Initialize customer service
	sc, err := Customer.New(customerID)
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// Get customer information
	customer, err := sc.Info().Get()
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// Bind data to struct
	if err = c.Bind(customer); err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// Validate struct
	if err = customer.Validate(); err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// Update information
	if err = sc.Info().Update(customer); err != nil {
		return SendErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	// Send response
	return SendOKResponse(c, Customer.MsgInfoUpdated)
}
