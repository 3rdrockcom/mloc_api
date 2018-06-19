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
	FirstName    null.String `form:"R1" json:"R1"`
	MiddleName   null.String `form:"R2" json:"R2"`
	LastName     null.String `form:"R3" json:"R3"`
	Suffix       null.String `form:"R4" json:"R4"`
	Birthday     null.String `form:"R5" json:"R5"`
	Address1     null.String `form:"R6" json:"R6"`
	Address2     null.String `form:"R7" json:"R7"`
	Country      null.Int    `form:"R8" json:"R8"`
	State        null.Int    `form:"R9" json:"R9"`
	City         null.Int    `form:"R10" json:"R10"`
	ZipCode      null.String `form:"R11" json:"R11"`
	HomeNumber   null.String `form:"R12" json:"R12"`
	MobileNumber null.String `form:"R13" json:"R13"`
	Email        null.String `form:"R14" json:"R14"`
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
	formKeys := []string{"R1", "R2", "R3", "R4", "R5", "R6", "R7", "R8", "R9", "R10", "R11", "R12", "R13", "R14"}
	for _, formKey := range formKeys {
		field := ""

		switch formKey {
		case "R1":
			field = "FirstName"
			customerBasic.FirstName = cr.FirstName
		case "R2":
			field = "MiddleName"
			customerBasic.MiddleName = cr.MiddleName
		case "R3":
			field = "LastName"
			customerBasic.LastName = cr.LastName
		case "R4":
			field = "Suffix"
			customerBasic.Suffix = cr.Suffix
		case "R5":
			field = "Birthday"
			customerBasic.Birthday = cr.Birthday
		case "R6":
			field = "Address1"
			customerBasic.Address1 = cr.Address1
		case "R7":
			field = "Address2"
			customerBasic.Address2 = cr.Address2
		case "R8":
			field = "Country"
			customerBasic.Country = cr.Country
		case "R9":
			field = "State"
			customerBasic.State = cr.State
		case "R10":
			field = "State"
			customerBasic.City = cr.City
		case "R11":
			field = "ZipCode"
			customerBasic.ZipCode = cr.ZipCode
		case "R12":
			field = "HomeNumber"
			customerBasic.HomeNumber = cr.HomeNumber
		case "R13":
			field = "MobileNumber"
			customerBasic.MobileNumber = cr.MobileNumber
		case "R14":
			field = "Email"
			customerBasic.Email = cr.Email
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
	CompanyName      null.String `form:"R1" json:"R1"`
	PhoneNumber      null.String `form:"R2" json:"R2"`
	NetPayPerCheck   null.Float  `form:"R3" json:"R3"`
	IncomeSource     null.Int    `form:"R4" json:"R4"`
	PayFrequency     null.Int    `form:"R5" json:"R5"`
	NextPayDate      null.String `form:"R6" json:"R6"`
	FollowingPayDate null.String `form:"R7" json:"R7"`
}

// Validate checks postform required is validation
func (ca CustomerAdditionalRequest) Validate() error {
	return validation.ValidateStruct(&ca,
		validation.Field(&ca.CompanyName, validation.Required),
		validation.Field(&ca.NetPayPerCheck, validation.Required, validation.Min(0.00)),
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
	formKeys := []string{"R1", "R2", "R3", "R4", "R5", "R6", "R7"}
	for _, formKey := range formKeys {
		field := ""

		switch formKey {
		case "R1":
			field = "CompanyName"
			customerAdditional.CompanyName = cr.CompanyName
		case "R2":
			field = "PhoneNumber"
			customerAdditional.PhoneNumber = cr.PhoneNumber
		case "R3":
			field = "NetPayPerCheck"
			customerAdditional.NetPayPerCheck = cr.NetPayPerCheck
		case "R4":
			field = "IncomeSource"
			customerAdditional.IncomeSource = cr.IncomeSource
		case "R5":
			field = "PayFrequency"
			customerAdditional.PayFrequency = cr.PayFrequency
		case "R6":
			field = "NextPayDate"
			customerAdditional.NextPayDate = cr.NextPayDate
		case "R7":
			field = "FollowingPayDate"
			customerAdditional.FollowingPayDate = cr.FollowingPayDate
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
	LoanAmount null.Float `form:"R2" json:"R2"`
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
	computedLoan, err := sc.Loan().ComputeLoanApplication(clr.LoanAmount.Float64)
	if err != nil {
		return err
	}

	// Send response
	return SendResponse(c, http.StatusOK, computedLoan)
}

// PostLoanApplicationRequest contains information for a loan application
type PostLoanApplicationRequest struct {
	LoanAmount null.Float `form:"R2" json:"R2"`
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
	err = sc.Loan().ProcessLoanApplication(lar.LoanAmount.Float64)
	if err != nil {
		return err
	}

	// Send response
	msg := Customer.MsgCustomerAppliedForLoan
	return SendOKResponse(c, msg)
}

// PostPayLoanRequest contains information about a loan payment
type PostPayLoanRequest struct {
	LoanAmount null.Float `form:"R2" json:"R2"`
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
	err = sc.Loan().ProcessLoanPayment(plr.LoanAmount.Float64)
	if err != nil {
		return err
	}

	// Send response
	msg := Customer.MsgCustomerAppliedForLoan
	return SendOKResponse(c, msg)
}
