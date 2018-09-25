package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/epointpayment/mloc_api_go/app/helpers"
	"github.com/epointpayment/mloc_api_go/app/models"
	Customer "github.com/epointpayment/mloc_api_go/app/services/customer"
	"github.com/epointpayment/mloc_api_go/app/services/payments"
	"github.com/epointpayment/mloc_api_go/app/services/payments/collection"
	"github.com/epointpayment/mloc_api_go/app/services/payments/registration"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/jinzhu/now"
	"github.com/labstack/echo"
	"github.com/shopspring/decimal"
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

	if customerInfo.NetPayPerCheck.Valid {
		customerInfo.NetPayPerCheckDisplay = null.StringFrom(decimal.NewFromFloat(customerInfo.NetPayPerCheck.Float64).StringFixed(helpers.DefaultCurrencyPrecision))
	}

	if customerInfo.CreditLimit.Valid {
		customerInfo.CreditLimitDisplay = null.StringFrom(decimal.NewFromFloat(customerInfo.CreditLimit.Float64).StringFixed(helpers.DefaultCurrencyPrecision))
	}

	if customerInfo.AvailableCredit.Valid {
		customerInfo.AvailableCreditDisplay = null.StringFrom(decimal.NewFromFloat(customerInfo.AvailableCredit.Float64).StringFixed(helpers.DefaultCurrencyPrecision))

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
	State        null.String `form:"R9" json:"R9"`
	City         null.String `form:"R10" json:"R10"`
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
			if val, err := strconv.Atoi(cr.State.String); err == nil {
				customerBasic.StateID = null.IntFrom(int64(val))
				customerBasic.State = null.NewString("", false)
			} else {
				customerBasic.StateID = null.NewInt(0, false)
				customerBasic.State = cr.State
			}
		case "R10":
			field = "City"
			if val, err := strconv.Atoi(cr.City.String); err == nil {
				customerBasic.CityID = null.IntFrom(int64(val))
				customerBasic.City = null.NewString("", false)
			} else {
				customerBasic.CityID = null.NewInt(0, false)
				customerBasic.City = cr.City
			}
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

			if field == "City" || field == "State" {
				fields = append(fields, field+"ID")
			}
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
	NetPayPerCheck   json.Number `form:"R3" json:"R3"`
	IncomeSource     null.Int    `form:"R4" json:"R4"`
	PayFrequency     null.Int    `form:"R5" json:"R5"`
	NextPayDate      null.String `form:"R6" json:"R6"`
	FollowingPayDate null.String `form:"R7" json:"R7"`
}

// Validate checks postform required is validation
func (ca CustomerAdditionalRequest) Validate() error {
	err := validation.ValidateStruct(&ca,
		validation.Field(&ca.CompanyName, validation.Required),
		validation.Field(&ca.NetPayPerCheck, validation.Required, validation.By(helpers.ValidateCurrencyAmount)),
		validation.Field(&ca.IncomeSource, validation.Required),
		validation.Field(&ca.PayFrequency, validation.Required),
		validation.Field(&ca.NextPayDate, validation.Date("2006-01-02")),
		validation.Field(&ca.FollowingPayDate, validation.Date("2006-01-02")),
	)
	if err != nil {
		return err
	}

	// Get current time
	t := time.Now().UTC()
	beginningOfYesterday := now.New(t.AddDate(0, 0, -1)).BeginningOfDay()
	var nextPayDate, followingPayDate time.Time

	// Validate next pay date
	if ca.NextPayDate.Valid {
		nextPayDate, err = time.Parse("2006-01-02", ca.NextPayDate.String)
		if err != nil {
			return err
		}

		// Next pay date must be equal or after yesterday
		if nextPayDate.Before(beginningOfYesterday) {
			return Customer.ErrInvalidNextPayDate
		}
	}

	// Validate following pay date
	if ca.FollowingPayDate.Valid {
		// Following pay date requires next pay date
		if !ca.NextPayDate.Valid {
			return Customer.ErrInvalidNextPayDate
		}

		followingPayDate, err = time.Parse("2006-01-02", ca.FollowingPayDate.String)
		if err != nil {
			return err
		}

		// Following pay date must be after next pay date
		if followingPayDate.Before(nextPayDate) || followingPayDate.Equal(nextPayDate) {
			return Customer.ErrInvalidFollowingPayDate
		}
	}

	return nil
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
		switch err {
		case Customer.ErrInvalidNextPayDate, Customer.ErrInvalidFollowingPayDate:
			break
		default:
			err = Customer.ErrCustomerIncompleteInfo
		}

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
			netPayPerCheckDecimal, _ := decimal.NewFromString(cr.NetPayPerCheck.String())
			netPayPerCheck, _ := netPayPerCheckDecimal.Float64()
			customerAdditional.NetPayPerCheck = null.FloatFrom(netPayPerCheck)
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
	CustomerID null.Int    `json:"fk_customer_id"`
	Mode       null.String `json:"mode"`
	Amount     string      `json:"amount"`
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
			Mode:       transactionHistory[i].Mode,
			Amount: decimal.NewFromFloat(transactionHistory[i].Amount.Float64).
				StringFixed(helpers.DefaultCurrencyPrecision),
			Type: transactionHistory[i].Type,
			Date: date,
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

	if customerLoanTotal.TotalPrincipalAmount.Valid {
		customerLoanTotal.TotalPrincipalAmountDisplay = null.StringFrom(decimal.NewFromFloat(customerLoanTotal.TotalPrincipalAmount.Float64).StringFixed(helpers.DefaultCurrencyPrecision))
	}

	if customerLoanTotal.TotalFeeAmount.Valid {
		customerLoanTotal.TotalFeeAmountDisplay = null.StringFrom(decimal.NewFromFloat(customerLoanTotal.TotalFeeAmount.Float64).StringFixed(helpers.DefaultCurrencyPrecision))
	}

	if customerLoanTotal.TotalAmount.Valid {
		customerLoanTotal.TotalAmountDisplay = null.StringFrom(decimal.NewFromFloat(customerLoanTotal.TotalAmount.Float64).StringFixed(helpers.DefaultCurrencyPrecision))
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

	// Get detailed customer information
	customerInfo, err := sc.Info().GetDetails()
	if err != nil {
		return err
	}

	// Register identifiers
	// Initialize payment service
	ps := payments.New()
	if err != nil {
		return err
	}

	if customerInfo.CLABE.String == "" {

		// Prepare registration request
		registrationRequest := registration.Request{
			Method:   payments.MethodSTP,
			Customer: *customerInfo,
		}

		// Execute payment registration
		registrationResponse, err := ps.Register(registrationRequest)
		if err != nil {
			// err = ErrIssuerFailedTransfer
			return err
		}

		customerBasic := new(models.CustomerBasic)
		customerBasic.CLABE = null.NewString(registrationResponse.Identifier, true)

		err = sc.Info().UpdateCustomerBasic(customerBasic, "CLABE")
		if err != nil {
			return err
		}
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
	LoanAmount null.String `form:"R2" json:"R2"`
}

// ComputeLoanResponse contains information about a loan
type ComputeLoanResponse struct {
	AvailableCredit  string `json:"available_credit"`
	Amount           string `json:"amount"`
	Fee              string `json:"fee"`
	Interest         string `json:"interest"`
	DateApplied      string `json:"date_applied"`
	DueDate          string `json:"due_date"`
	DueDateFormatted string `json:"due_date_formatted"`
	TotalAmount      string `json:"total_amount"`
}

// Validate checks postform required is validation
func (clr PostComputeLoanRequest) Validate() error {
	return validation.ValidateStruct(&clr,
		validation.Field(&clr.LoanAmount, validation.Required, validation.By(helpers.ValidateCurrencyAmount)),
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

	// Convert loan amount to decimal
	loanAmount, _ := decimal.NewFromString(clr.LoanAmount.String)
	loanAmount = loanAmount.RoundBank(helpers.DefaultCurrencyPrecision)

	// Calculate loan application
	computedLoan, err := sc.Loan().ComputeLoanApplication(loanAmount)
	if err != nil {
		return err
	}

	computedLoanResponse := ComputeLoanResponse{
		AvailableCredit:  computedLoan.AvailableCredit.StringFixed(helpers.DefaultCurrencyPrecision),
		Amount:           computedLoan.Amount.StringFixed(helpers.DefaultCurrencyPrecision),
		Fee:              computedLoan.Fee.StringFixed(helpers.DefaultCurrencyPrecision),
		Interest:         computedLoan.Interest.StringFixed(helpers.DefaultCurrencyPrecision),
		DateApplied:      computedLoan.DateApplied,
		DueDate:          computedLoan.DueDate,
		DueDateFormatted: computedLoan.DueDateFormatted,
		TotalAmount:      computedLoan.TotalAmount.StringFixed(helpers.DefaultCurrencyPrecision),
	}

	// Send response
	return SendResponse(c, http.StatusOK, computedLoanResponse)
}

// PostLoanApplicationRequest contains information for a loan application
type PostLoanApplicationRequest struct {
	DisbursementMethod null.String `form:"R3" json:"R3"`
	BankAccountID      null.Int    `form:"R4" json:"R4"`
	LoanAmount         null.String `form:"R2" json:"R2"`
}

// Validate checks postform required is validation
func (lar PostLoanApplicationRequest) Validate() error {
	return validation.ValidateStruct(&lar,
		validation.Field(&lar.LoanAmount, validation.Required, validation.By(helpers.ValidateCurrencyAmount)),
		validation.Field(&lar.DisbursementMethod, validation.Required, validation.In(payments.MethodEPOINT, payments.MethodSTP)),
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

	// Get disbursement information
	disbursementMethod := lar.DisbursementMethod.String
	disbursementBankAccount := int(lar.BankAccountID.Int64)

	// Convert loan amount to decimal
	loanAmount, _ := decimal.NewFromString(lar.LoanAmount.String)
	loanAmount = loanAmount.RoundBank(helpers.DefaultCurrencyPrecision)

	// Calculate loan application
	err = sc.Loan().ProcessLoanApplication(disbursementMethod, disbursementBankAccount, loanAmount)
	if err != nil {
		return err
	}

	// Send response
	msg := Customer.MsgCustomerAppliedForLoan
	return SendOKResponse(c, msg)
}

// PostPayLoanRequest contains information about a loan payment
type PostPayLoanRequest struct {
	LoanAmount null.String `form:"R2" json:"R2"`
}

// Validate checks postform required is validation
func (plr PostPayLoanRequest) Validate() error {
	return validation.ValidateStruct(&plr,
		validation.Field(&plr.LoanAmount, validation.Required, validation.By(helpers.ValidateCurrencyAmount)),
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

	// Convert payment amount to decimal
	loanAmount, _ := decimal.NewFromString(plr.LoanAmount.String)
	loanAmount = loanAmount.RoundBank(helpers.DefaultCurrencyPrecision)

	// Calculate loan payment
	err = sc.Loan().DoCollection(loanAmount)
	if err != nil {
		return err
	}

	// Send response
	msg := Customer.MsgCustomerMadeLoanPayment
	return SendOKResponse(c, msg)
}

// PostPayLoanRequest contains information about a loan payment
type PostPayLoanPushRequest struct {
	CollectionMethod       null.String `form:"R3" json:"R3"`
	TransactionAmount      null.String `form:"R2" json:"R2"`
	TransactionDate        null.String `form:"R4" json:"R4"`
	TransactionID          null.String `form:"R5" json:"R5"`
	TransactionDescription null.String `form:"R6" json:"R6"`
}

// Validate checks postform required is validation
func (plr PostPayLoanPushRequest) Validate() error {
	return validation.ValidateStruct(&plr,
		validation.Field(&plr.TransactionAmount, validation.Required, validation.By(helpers.ValidateCurrencyAmount)),
		validation.Field(&plr.CollectionMethod, validation.Required, validation.In(payments.MethodSTP)),
		validation.Field(&plr.TransactionDate, validation.Required, validation.Date("2006-01-02 15:04:05")),
		validation.Field(&plr.TransactionID, validation.Required),
	)
}

// PostPayLoan processes a loan payment
func (co *Controllers) PostPayLoanPush(c echo.Context) error {
	// Get customer ID
	customerID := c.Get("customerID").(int)

	// Initialize customer service
	sc, err := Customer.New(customerID)
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	plr := PostPayLoanPushRequest{}

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

	// Convert payment amount to decimal
	loanAmount, _ := decimal.NewFromString(plr.TransactionAmount.String)
	loanAmount = loanAmount.RoundBank(helpers.DefaultCurrencyPrecision)

	// Date
	t, err := time.Parse("2006-01-02 15:04:05", plr.TransactionDate.String)
	if err != nil {
		return err
	}

	collectionPush := collection.Push{
		Method: plr.CollectionMethod.String,
		Transaction: collection.PushTransaction{
			ID:          plr.TransactionID.String,
			Date:        t,
			Amount:      loanAmount,
			Description: plr.TransactionDescription.String,
		},
	}

	// Calculate loan payment
	err = sc.Loan().DoCollectionPush(collectionPush)
	if err != nil {
		return err
	}

	// Send response
	msg := Customer.MsgCustomerMadeLoanPayment
	return SendOKResponse(c, msg)
}

// GetBankAccount gets information about a customer's bank account(s)
func (co *Controllers) GetBankAccount(c echo.Context) (err error) {
	// Get customer ID
	customerID := c.Get("customerID").(int)

	// Get bank account id
	queryBankAccountID := c.QueryParam("R2")
	bankAccountID, err := strconv.Atoi(queryBankAccountID)
	if err != nil && queryBankAccountID != "" {
		err = Customer.ErrInvalidBankAccountID
		return
	}

	// Initialize customer service
	sc, err := Customer.New(customerID)
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// Get bank account
	if bankAccountID > 0 {
		bankAccount, err := sc.BankAccount().Get(bankAccountID)
		if err != nil {
			return err
		}

		bankAccount.DateCreatedDisplay = bankAccount.DateCreated.Format("2006-01-02 15:04:05")

		if !bankAccount.DateUpdated.IsZero() {
			bankAccount.DateUpdatedDisplay = bankAccount.DateUpdated.Format("2006-01-02 15:04:05")
		}

		// Send response
		return SendResponse(c, http.StatusOK, bankAccount)
	}

	// Get list of bank accounts
	bankAccounts, err := sc.BankAccount().GetAll()
	if err != nil {
		return err
	}

	for i := range bankAccounts {
		bankAccounts[i].DateCreatedDisplay = bankAccounts[i].DateCreated.Format("2006-01-02 15:04:05")
		if !bankAccounts[i].DateUpdated.IsZero() {
			bankAccounts[i].DateUpdatedDisplay = bankAccounts[i].DateUpdated.Format("2006-01-02 15:04:05")
		}

	}

	// Send response
	return SendResponse(c, http.StatusOK, bankAccounts)
}

// PostCreateBankAccountRequest contains information about a bank account creation request
type PostCreateBankAccountRequest struct {
	Alias         null.String `form:"R2" json:"R2"`
	AccountNumber null.String `form:"R3" json:"R3"`
	BankCode      null.String `form:"R4" json:"R4"`
}

// Validate checks postform required is validation
func (r PostCreateBankAccountRequest) Validate() error {
	return validation.ValidateStruct(&r,
		// validation.Field(&r.Alias, validation.Required),
		validation.Field(&r.AccountNumber, validation.Required),
	)
}

// PostCreateBankAccount creates a bank account for a customer
func (co *Controllers) PostCreateBankAccount(c echo.Context) (err error) {
	// Get customer ID
	customerID := c.Get("customerID").(int)

	r := PostCreateBankAccountRequest{}

	// Bind data to struct
	if err = c.Bind(&r); err != nil {
		err = Customer.ErrInvalidBankAccountInformation
		return err
	}

	// Validate struct
	if err = r.Validate(); err != nil {
		err = Customer.ErrInvalidBankAccountInformation
		return err
	}

	// Initialize customer service
	sc, err := Customer.New(customerID)
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// Create bank account
	bankAccount := models.CustomerBankAccount{
		Alias:         r.Alias.String,
		AccountNumber: r.AccountNumber.String,
		BankCode:      r.BankCode.String,
	}
	err = sc.BankAccount().Create(&bankAccount)
	if err != nil {
		return err
	}

	bankAccount.DateCreatedDisplay = bankAccount.DateCreated.Format("2006-01-02 15:04:05")

	// Send response
	return SendResponse(c, http.StatusOK, bankAccount)
}

// PostUpdateBankAccountRequest contains information about a bank account update request
type PostUpdateBankAccountRequest struct {
	BankAccountID null.Int    `form:"R2" json:"R2"`
	Alias         null.String `form:"R3" json:"R3"`
}

// Validate checks postform required is validation
func (r PostUpdateBankAccountRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.BankAccountID, validation.Required, validation.Min(0)),
		validation.Field(&r.Alias, validation.Required),
	)
}

// PostUpdateBankAccount updates a customer's bank account information
func (co *Controllers) PostUpdateBankAccount(c echo.Context) (err error) {
	// Get customer ID
	customerID := c.Get("customerID").(int)

	r := PostUpdateBankAccountRequest{}

	// Bind data to struct
	if err = c.Bind(&r); err != nil {
		err = Customer.ErrInvalidBankAccountInformation
		return err
	}

	// Validate struct
	if err = r.Validate(); err != nil {
		err = Customer.ErrInvalidBankAccountInformation
		return err
	}

	// Initialize customer service
	sc, err := Customer.New(customerID)
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// Get bank account
	bankAccountID := int(r.BankAccountID.Int64)
	bankAccount, err := sc.BankAccount().Get(bankAccountID)
	if err != nil {
		return err
	}

	// Prepare data
	bankAccount.Alias = r.Alias.String

	// Update bank account
	err = sc.BankAccount().Update(bankAccount)
	if err != nil {
		return
	}

	bankAccount.DateCreatedDisplay = bankAccount.DateCreated.Format("2006-01-02 15:04:05")
	bankAccount.DateUpdatedDisplay = bankAccount.DateUpdated.Format("2006-01-02 15:04:05")

	// Send response
	return SendOKResponse(c, Customer.MsgCustomerUpdatedCustomerBankAccount)
}

// PostDeleteBankAccountRequest contains information about a bank account deletion request
type PostDeleteBankAccountRequest struct {
	BankAccountID null.Int `form:"R2" json:"R2"`
}

// Validate checks postform required is validation
func (r PostDeleteBankAccountRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.BankAccountID, validation.Required, validation.Min(0)),
	)
}

// PostDeleteBankAccount deletes a customer's bank account
func (co *Controllers) PostDeleteBankAccount(c echo.Context) (err error) {
	// Get customer ID
	customerID := c.Get("customerID").(int)

	r := PostDeleteBankAccountRequest{}

	// Bind data to struct
	if err = c.Bind(&r); err != nil {
		err = Customer.ErrInvalidBankAccountInformation
		return err
	}

	// Validate struct
	if err = r.Validate(); err != nil {
		err = Customer.ErrInvalidBankAccountInformation
		return err
	}

	// Initialize customer service
	sc, err := Customer.New(customerID)
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// Get bank account
	bankAccountID := int(r.BankAccountID.Int64)
	bankAccount, err := sc.BankAccount().Get(bankAccountID)
	if err != nil {
		return err
	}

	//
	err = sc.BankAccount().Delete(bankAccount)
	if err != nil {
		return err
	}

	// Send response
	return SendOKResponse(c, Customer.MsgCustomerDeletedCustomerBankAccount)
}
