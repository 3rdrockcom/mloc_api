package controllers

import (
	"net/http"
	"strings"
	"time"

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

	if customerInfo.BirthDate.Valid {
		t, err := time.Parse(time.RFC3339, customerInfo.BirthDate.String)
		if err != nil {
			return err
		}
		customerInfo.BirthDate = null.StringFrom(t.Format("2006-01-02"))
	}

	if customerInfo.NextPayDate.Valid {
		t, err := time.Parse(time.RFC3339, customerInfo.NextPayDate.String)
		if err != nil {
			return err
		}
		customerInfo.NextPayDate = null.StringFrom(t.Format("2006-01-02"))
	}

	return SendResponse(c, http.StatusOK, customerInfo)
}

// CustomerBasicRequest contains basic information of a customer
type CustomerBasicRequest struct {
	FirstName    string `form:"R1" json:"R1"`
	MiddleName   string `form:"R2" json:"R2"`
	LastName     string `form:"R3" json:"R3"`
	Suffix       string `form:"R4" json:"R4"`
	Birthday     string `form:"R5" json:"R5"`
	Address1     string `form:"R6" json:"R6"`
	Address2     string `form:"R7" json:"R7"`
	Country      int64  `form:"R8" json:"R8"`
	State        int64  `form:"R9" json:"R9"`
	City         int64  `form:"R10" json:"R10"`
	ZipCode      string `form:"R11" json:"R11"`
	HomeNumber   string `form:"R12" json:"R12"`
	MobileNumber string `form:"R13" json:"R13"`
	Email        string `form:"R14" json:"R14"`
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

// PostCustomerBasic updates customer basic information
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

// CustomerAdditionalRequest contains additional information of a customer
type CustomerAdditionalRequest struct {
	CompanyName      string  `form:"R1" json:"R1"`
	PhoneNumber      string  `form:"R2" json:"R2"`
	NetPayPerCheck   float64 `form:"R3" json:"R3"`
	IncomeSource     int64   `form:"R4" json:"R4"`
	PayFrequency     int64   `form:"R5" json:"R5"`
	NextPayDate      string  `form:"R6" json:"R6"`
	FollowingPayDate string  `form:"R7" json:"R7"`
}

// Validate checks postform required is validation
func (ca CustomerAdditionalRequest) Validate() error {
	return validation.ValidateStruct(&ca,
		validation.Field(&ca.CompanyName, validation.Required),
		validation.Field(&ca.NetPayPerCheck, validation.Required),
		validation.Field(&ca.IncomeSource, validation.Required),
		validation.Field(&ca.PayFrequency, validation.Required),
		validation.Field(&ca.NextPayDate, validation.Date("2006-01-02")),
		validation.Field(&ca.FollowingPayDate, validation.Date("2006-01-02")),
	)
}

// PostCustomerAdditional updates customer additional information
func (co Controllers) PostCustomerAdditional(c echo.Context) error {
	// Get customer ID
	customerID := c.Get("customerID").(int)

	// Initialize customer service
	sc, err := Customer.New(customerID)
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	cr := CustomerAdditionalRequest{}

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
	customerAdditional := new(models.CustomerAdditional)
	customerAdditional.ID = customerID

	fields := []string{}
	formKeys, _ := c.FormParams() // previously called in form binding, no need to check errors
	for formKey := range formKeys {
		field := ""

		switch formKey {
		case "R1":
			field = "CompanyName"
			customerAdditional.CompanyName = null.StringFrom(cr.CompanyName)
		case "R2":
			field = "PhoneNumber"
			customerAdditional.PhoneNumber = null.StringFrom(cr.PhoneNumber)
		case "R3":
			field = "NetPayPerCheck"
			customerAdditional.NetPayPerCheck = null.FloatFrom(cr.NetPayPerCheck)
		case "R4":
			field = "IncomeSource"
			customerAdditional.IncomeSource = null.IntFrom(cr.IncomeSource)
		case "R5":
			field = "PayFrequency"
			customerAdditional.PayFrequency = null.IntFrom(cr.PayFrequency)
		case "R6":
			field = "NextPayDate"
			customerAdditional.NextPayDate = null.StringFrom(cr.NextPayDate)
		case "R7":
			field = "FollowingPayDate"
			customerAdditional.FollowingPayDate = null.StringFrom(cr.FollowingPayDate)
		}

		if field != "" {
			fields = append(fields, field)
		}
	}

	// Update information
	if err = sc.Info().UpdateCustomerAdditional(customerAdditional, fields...); err != nil {
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
	switch strings.ToUpper(transactionType) {
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

// PostAcceptTermsAndConditions is called when a customer has accepted the terms and conditions
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

// PostCreditLineApplication processes a credit line application
func (co *Controllers) PostCreditLineApplication(c echo.Context) error {
	// Get customer ID
	customerID := c.Get("customerID").(int)

	// Initialize customer service
	sc, err := Customer.New(customerID)
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// Get customer loan total
	_, err = sc.Loan().ProcessCreditLineApplication()
	if err != nil {
		return err
	}

	// Send response
	msg := Customer.MsgCustomerAddedCreditLine
	return SendOKResponse(c, msg)
}

// PostComputeLoanRequest contains information for computing a loan application
type PostComputeLoanRequest struct {
	LoanAmount float64 `form:"R2" json:"R2"`
}

// Validate checks postform required is validation
func (clr PostComputeLoanRequest) Validate() error {
	return validation.ValidateStruct(&clr,
		validation.Field(&clr.LoanAmount, validation.Required, validation.Min(0.00)),
	)
}

// PostComputeLoan computes a loan application
func (co *Controllers) PostComputeLoan(c echo.Context) error {
	// Get customer ID
	customerID := c.Get("customerID").(int)

	// Initialize customer service
	sc, err := Customer.New(customerID)
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	clr := PostComputeLoanRequest{}

	// Bind data to struct
	if err = c.Bind(&clr); err != nil {
		err = Customer.ErrInvalidLoanAmount
		return err
	}

	// Validate struct
	if err = clr.Validate(); err != nil {
		err = Customer.ErrInvalidLoanAmount
		return err
	}

	// Calculate loan application
	computedLoan, err := sc.Loan().ComputeLoanApplication(clr.LoanAmount)
	if err != nil {
		return err
	}

	// Send response
	return SendResponse(c, http.StatusOK, computedLoan)
}

// PostLoanApplicationRequest contains information for a loan application
type PostLoanApplicationRequest struct {
	LoanAmount float64 `form:"R2" json:"R2"`
}

// Validate checks postform required is validation
func (lar PostLoanApplicationRequest) Validate() error {
	return validation.ValidateStruct(&lar,
		validation.Field(&lar.LoanAmount, validation.Required, validation.Min(0.00)),
	)
}

// PostLoanApplication processes a loan application
func (co *Controllers) PostLoanApplication(c echo.Context) error {
	// Get customer ID
	customerID := c.Get("customerID").(int)

	// Initialize customer service
	sc, err := Customer.New(customerID)
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	lar := PostLoanApplicationRequest{}

	// Bind data to struct
	if err = c.Bind(&lar); err != nil {
		err = Customer.ErrInvalidLoanAmount
		return err
	}

	// Validate struct
	if err = lar.Validate(); err != nil {
		err = Customer.ErrInvalidLoanAmount
		return err
	}

	// Calculate loan application
	err = sc.Loan().ProcessLoanApplication(lar.LoanAmount)
	if err != nil {
		return err
	}

	// Send response
	msg := Customer.MsgCustomerAppliedForLoan
	return SendOKResponse(c, msg)
}

// PostPayLoanRequest contains information about a loan payment
type PostPayLoanRequest struct {
	LoanAmount float64 `form:"R2" json:"R2"`
}

// Validate checks postform required is validation
func (plr PostPayLoanRequest) Validate() error {
	return validation.ValidateStruct(&plr,
		validation.Field(&plr.LoanAmount, validation.Required, validation.Min(0.00)),
	)
}

// PostPayLoan processes a loan payment
func (co *Controllers) PostPayLoan(c echo.Context) error {
	// Get customer ID
	customerID := c.Get("customerID").(int)

	// Initialize customer service
	sc, err := Customer.New(customerID)
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	plr := PostPayLoanRequest{}

	// Bind data to struct
	if err = c.Bind(&plr); err != nil {
		err = Customer.ErrInvalidLoanAmount
		return err
	}

	// Validate struct
	if err = plr.Validate(); err != nil {
		err = Customer.ErrInvalidLoanAmount
		return err
	}

	// Calculate loan payment
	err = sc.Loan().ProcessLoanPayment(plr.LoanAmount)
	if err != nil {
		return err
	}

	// Send response
	msg := Customer.MsgCustomerAppliedForLoan
	return SendOKResponse(c, msg)
}
