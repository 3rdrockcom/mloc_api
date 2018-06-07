package controllers

import (
	"net/http"

	"github.com/epointpayment/mloc_api_go/app/models"

	Customer "github.com/epointpayment/mloc_api_go/app/services/customer"
	"github.com/labstack/echo"
	null "gopkg.in/guregu/null.v3"
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

// TransactionsHistoryResponse is array of transactions
type TransactionsHistoryResponse []TransactionResponse

// TransactionResponse contains information about a loan transaction
type TransactionResponse struct {
	CustomerID null.Int    ` json:"fk_customer_id"`
	Amount     null.Float  `json:"amount"`
	Type       string      `json:"t_type"`
	Date       null.String `json:"t_date"`
}

// GetTransactionHistory displays transaction history of a customer based on transaction type
func (co *Controllers) GetTransactionHistory(c echo.Context) error {
	// Get customer ID
	customerID := c.Get("customerID").(int)

	// Get transaction type
	transactionType := c.QueryParam("R2")

	// Initialize customer service
	sc, err := Customer.New(customerID)
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// Get transaction history
	transactionHistory := models.TransactionsHistory{}
	switch transactionType {
	case "LOAN":
		fallthrough
	case "PAYMENT":
		transactionHistory, err = sc.Loan().GetTransactionHistoryByType(transactionType)
		if err != nil {
			return err
		}
	default:
		transactionHistory, err = sc.Loan().GetTransactionHistory()
		if err != nil {
			return err
		}
	}

	// Convert to response format
	transactionHistoryResponse := TransactionsHistoryResponse{}
	for i := range transactionHistory {
		// Convert to specific date format
		date := null.NewString("", false)
		if transactionHistory[i].Date.Valid {
			date = null.StringFrom(transactionHistory[i].Date.Ptr().Format("2006-01-02 15:04:05"))
		}

		transactionHistoryResponse = append(transactionHistoryResponse, TransactionResponse{
			CustomerID: transactionHistory[i].CustomerID,
			Amount:     transactionHistory[i].Amount,
			Type:       transactionHistory[i].Type,
			Date:       date,
		})
	}

	// Send response
	return SendResponse(c, http.StatusOK, transactionHistoryResponse)
}

// GetCustomerLoan displays information about a customer's loan total
func (co *Controllers) GetCustomerLoan(c echo.Context) error {
	// Get customer ID
	customerID := c.Get("customerID").(int)

	// Initialize customer service
	sc, err := Customer.New(customerID)
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// Get customer loan total
	customerLoanTotal, err := sc.Loan().GetCustomerLoanTotal()
	if err != nil {
		return err
	}

	// Send response
	return SendResponse(c, http.StatusOK, customerLoanTotal)
}

// PostAcceptTermsAndCondition accepts the term and condition
// the value of accept will store in tblcustomeragreement
func (co *Controllers) PostAcceptTermsAndConditions(c echo.Context) error {
	// Get customer ID
	customerID := c.Get("customerID").(int)

	// Initialize customer service
	sc, err := Customer.New(customerID)
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// Get customer loan total
	err = sc.Loan().AcceptedCustomerAgreement()
	if err != nil {
		return err
	}

	// Send response
	msg := Customer.MsgCustomerAcceptedAgreement
	return SendOKResponse(c, msg)

}
