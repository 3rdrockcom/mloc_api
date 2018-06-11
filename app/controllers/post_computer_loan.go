package controllers

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/labstack/echo"
)

// GetCustomer contains informatin of vew_customer_info
type GetCustomer struct {
	CustomerID        int
	FirstName         *string
	AvailableCredit   *float64
	CustUniqueID      *string
	CreditLineID      *int
	ProgramCustomerID *int
	MobileNumber      *string
}

// Fee contains informatin of tblfee
type Fee struct {
	Active     *string
	Percentage *float64
	Fixed      *float64
}

// Interest contains informatin of tblinterest
type Interest struct {
	Active     *string
	Percentage *float64
	Fixed      *float64
}

// ComputedFee contains informatin of loan
type ComputedFee struct {
	AvailableCredit  *float64
	Amount           *float64
	Fee              float64
	Interest         float64
	DateApplied      string // may need to change date type
	DueDate          string // may need to change date type
	DueDateFormatted string // may need to change date type
	TotalAmount      float64
}

// TableName gets tblfee from database
func (a Fee) TableName() string {
	return "tblfee"
}

// TableName gets tblinterest from database
func (a Interest) TableName() string {
	return "tblinterest"
}

// TableName gets view_customer_info from database
func (a GetCustomer) TableName() string {
	return "view_customer_info"
}

var (

	// ErrNoValidAmount is given when the customer input amount is not valid
	ErrNoValidAmount = "amount is not a valid"

	// ErrNotEnoughCredit is given when amount is larger than credit
	ErrNotEnoughCredit = "You dont have enough available credit."
)

// PostComputerLoan proceeds to loan application of the customer
func (co *Controllers) PostComputerLoan(c echo.Context) error {

	// Get custoemrID
	customerID := c.Get("customerID").(int)
	// Get customer Unique ID and amount
	custUniqueID := c.FormValue("R1")
	amount := c.FormValue("R2")

	loanAmount, err := strconv.ParseFloat(amount, 64) // Convert amount string type to float64 type
	if err != nil {
		message := ErrNoValidAmount
		return SendErrorResponse(c, http.StatusBadRequest, message)
	}

	// call function to compute process amount and duedate
	update, err := ComputerLoanApplication(customerID, custUniqueID, loanAmount)
	//	fmt.Println(err)
	if err != nil {
		message := ErrNotEnoughCredit
		return SendErrorResponse(c, http.StatusBadRequest, message)
	}

	// Success to compter loan application
	return SendResponse(c, http.StatusOK, update)

}

// ComputerLoanApplication gets total amount loan of customer
func ComputerLoanApplication(custID int, custUniqueID string, amount float64) (ComputedFee, error) {
	computedfee := ComputedFee{}
	totalAmount := amount
	customer := GetCustomer{}
	var feeAmount float64
	var interestAmount float64

	// Get infomation of customer
	q := db.Select().
		From("view_customer_info").
		Where(dbx.HashExp{"customer_id": custID}) // cust_unique_id is *string type
	err := q.One(&customer)

	if err != nil {
		message := ErrNoCustomerExist
		return computedfee, errors.New(message)
	}

	if customer.CustomerID == 0 || amount > *customer.AvailableCredit {
		message := ErrNoCustomerExist
		return computedfee, errors.New(message)
	}

	// Compute fee
	fee, err := getFee()
	fmt.Println(fee)
	if err != nil {
		message := "can't get fee"
		return computedfee, errors.New(message)
	}

	if *fee.Percentage > 0.0 {
		feeAmount = (*fee.Percentage / 100.00) * amount
		feeAmount = math.Round(feeAmount*100) / 100

	} else {
		feeAmount = *fee.Fixed
	}
	fmt.Println(feeAmount)

	// Compute interest
	interest, err := getInterest()

	if *interest.Percentage > 0.0 {
		interestAmount = (*interest.Percentage / 100.00) * amount
		interestAmount = math.Round(interestAmount*100) / 100
	} else {
		interestAmount = *interest.Fixed
	}

	customerCreditLineID := strconv.Itoa(*customer.CreditLineID) // CreditLineID is convered to string for pass function
	tierNoDays, err := GetLoanCreditLimit(customerCreditLineID)  // it returns struct of creditlimit
	if err != nil {
		message := ErrNoCustomerExist
		return computedfee, errors.New(message)
	}

	dueDate := time.Now().AddDate(0, 0, *tierNoDays.NoOfDays).Format("2006-01-02 15:04:05")
	fmt.Println(dueDate)

	dueDateFormatted := time.Now().AddDate(0, 0, *tierNoDays.NoOfDays).Format("01-02-2006 03:04 PM")
	fmt.Println(dueDateFormatted)

	// Compute total Amount
	totalAmount += feeAmount
	totalAmount += interestAmount

	customerAvailableCredit := math.Round(*customer.AvailableCredit*100) / 100
	customer.AvailableCredit = &customerAvailableCredit

	amount = math.Round(amount*100) / 100
	totalAmount = math.Round(totalAmount*100) / 100

	computedfee.AvailableCredit = customer.AvailableCredit
	computedfee.Amount = &amount
	computedfee.Fee = feeAmount
	computedfee.Interest = interestAmount
	computedfee.DateApplied = time.Now().Format("2006-01-02 15:04:05")
	computedfee.DueDate = dueDate
	computedfee.DueDateFormatted = dueDateFormatted
	computedfee.TotalAmount = totalAmount
	return computedfee, nil // Return struct of computedfee

}

// getFee gets tblfee information from database and returns struct of Fee
func getFee() (Fee, error) {
	fee := Fee{}
	dbFee := db.Select().
		From("tblFee").
		Where(dbx.HashExp{"active": "YES"}) // cust_unique_id is *string type
	err := dbFee.One(&fee)
	if err != nil {
		message := "can't get fee"
		return fee, errors.New(message)
	}

	return fee, nil

}

// getInterest gets tblInterest from database and returns struct of Interest
func getInterest() (Interest, error) {
	interest := Interest{}
	dbinterest := db.Select().
		From("tblInterest").
		Where(dbx.HashExp{"active": "YES"}) // cust_unique_id is *string type
	err := dbinterest.One(&interest)
	if err != nil {
		message := "can't get interest"
		return interest, errors.New(message)
	}

	return interest, nil

}
