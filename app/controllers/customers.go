package controllers

import (
	"net/http"

	"github.com/epointpayment/mloc_api_go/app/models"
	Customer "github.com/epointpayment/mloc_api_go/app/services/customer"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
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

// CustomerBasicRequest contains basic information of a customer
type CustomerBasicRequest struct {
	ID           int
	FirstName    string `form:"R1"`
	MiddleName   string `form:"R2"`
	LastName     string `form:"R3"`
	Suffix       string `form:"R4"`
	Birthday     string `form:"R5"` // may need to change date type
	Address1     string `form:"R6"`
	Address2     string `form:"R7"`
	Country      int64  `form:"R8"`
	State        int64  `form:"R9"`
	City         int64  `form:"R10"`
	ZipCode      string `form:"R11"`
	HomeNumber   string `form:"R12"`
	MobileNumber string `form:"R13"`
	Email        string `form:"R14"`
}

// Validate checks postform required is validation
func (cb CustomerBasicRequest) Validate() error {
	return validation.ValidateStruct(&cb,
		validation.Field(&cb.FirstName, validation.Required),
		validation.Field(&cb.LastName, validation.Required),
		validation.Field(&cb.MobileNumber, validation.Required),
		validation.Field(&cb.Email, validation.Required, is.Email),
		validation.Field(&cb.Birthday, validation.Date("2006-01-02")),
	)
}

// PostCustomerBasic updates customer information
func (co Controllers) PostCustomerBasic(c echo.Context) error {
	// Get customer ID
	customerID := c.Get("customerID").(int)

	// Initialize customer service
	sc, err := Customer.New(customerID)
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	cr := CustomerBasicRequest{}

	// Bind data to struct
	if err = c.Bind(&cr); err != nil {
		err = Customer.ErrProblemOccured
		return err
	}

	// Validate struct
	if err = cr.Validate(); err != nil {
		err = Customer.ErrCustomerIncompleteInfo
		return err
	}

	// Prepare customer information
	customerBasic := new(models.CustomerBasic)
	customerBasic.ID = customerID

	fields := []string{}
	formKeys, _ := c.FormParams() // previously called in form binding, no need to check errors
	for formKey := range formKeys {
		field := ""

		switch formKey {
		case "R1":
			field = "FirstName"
			customerBasic.FirstName = null.StringFrom(cr.FirstName)
		case "R2":
			field = "MiddleName"
			customerBasic.MiddleName = null.StringFrom(cr.MiddleName)
		case "R3":
			field = "LastName"
			customerBasic.LastName = null.StringFrom(cr.LastName)
		case "R4":
			field = "Suffix"
			customerBasic.Suffix = null.StringFrom(cr.Suffix)
		case "R5":
			field = "Birthday"
			customerBasic.Birthday = null.StringFrom(cr.Birthday)
		case "R6":
			field = "Address1"
			customerBasic.Address1 = null.StringFrom(cr.Address1)
		case "R7":
			field = "Address2"
			customerBasic.Address2 = null.StringFrom(cr.Address2)
		case "R8":
			field = "Country"
			customerBasic.Country = null.IntFrom(cr.Country)
		case "R9":
			field = "State"
			customerBasic.State = null.IntFrom(cr.State)
		case "R10":
			field = "State"
			customerBasic.City = null.IntFrom(cr.City)
		case "R11":
			field = "ZipCode"
			customerBasic.ZipCode = null.StringFrom(cr.ZipCode)
		case "R12":
			field = "HomeNumber"
			customerBasic.HomeNumber = null.StringFrom(cr.HomeNumber)
		case "R13":
			field = "MobileNumber"
			customerBasic.MobileNumber = null.StringFrom(cr.MobileNumber)
		case "R14":
			field = "Email"
			customerBasic.Email = null.StringFrom(cr.Email)
		}

		if field != "" {
			fields = append(fields, field)
		}
	}

	// Update information
	if err = sc.Info().UpdateCustomerBasic(customerBasic, fields...); err != nil {
		err = Customer.ErrProblemOccured
		return err
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
