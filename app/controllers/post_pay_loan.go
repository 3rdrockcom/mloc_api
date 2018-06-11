package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/labstack/echo"
)

type CustomerLoans []CustomerLoan

type CustomerLoan struct {
	ID                  int
	FKCustomerID        *int
	LoanApplicationID   *int
	EpointTransactionID *string
	LoanIntervalID      *int
	LoanTermID          *int
	LoanAmount          *float64
	InterestAmount      *float64
	FeeAmount           *float64
	TotalAmount         *float64
	TotalPaidPrincipal  *float64
	TotalPaidFee        *float64
	TotalPaidAmount     *float64
	ReferenceCode       *string
	IsPaid              *int
	DueDate             *time.Time
	LoanDate            *time.Time
	ApprovedBy          *string
	ApprovedDate        *time.Time
	CreatedDate         time.Time
}

// CustomerPaymet contains information of customer payment
type CustomerPaymet struct {
	ID                  int
	FKCustomerID        *int
	ReferenceCode       *string
	EpointTransactionID *float64
	PayMentAmount       *float64
	DatePaid            *string // may need to change date type
	PaidBy              *string
}

// TableName gets tblcustoemrpayment from database
func (a CustomerPaymet) TableName() string {
	return "tblcustomerpayment"
}

// TableName gets tblcustomerloan from database
func (a CustomerLoan) TableName() string {
	return "tblcustomerloan"
}

var (
	// ErrNotEnoughWalletBalance is given if the customer doesn't have enought wallet balance
	ErrNotEnoughWalletBalance = "You dont have enough available wallet balance."

	// ErrProcessLoanPayment is given while process pay loan
	ErrProcessLoanPayment = "Error while processing loan payment."

	// MsgCustomerSuccessPaidLoan is given while customer paid loan is successsfully
	MsgCustomerSuccessPaidLoan = "Customer successfully paid a loan."
)

func (co *Controllers) postPayLoan(c echo.Context) error {
	// Get custoemrID
	customerID := c.Get("customerID").(int)
	// Get customer Unique ID and amount
	custUniqueID := c.FormValue("R1")
	amount := c.FormValue("R2")

	loanAmount, err := strconv.ParseFloat(amount, 64) // convert amount string type to float64 type
	if err != nil {
		message := ErrNoValidAmount
		return SendErrorResponse(c, http.StatusBadRequest, message)
	}

	insert := payLoanProcess(customerID, custUniqueID, loanAmount)

	if insert.Status == false {
		message := insert.StatusMsg
		return SendErrorResponse(c, http.StatusBadRequest, message)
	}

	// Loan application is success to apply
	message := insert.StatusMsg
	return SendOKResponse(c, message)
}

func payLoanProcess(custID int, custUniqueID string, amount float64) ReturnMsg {
	returnMsg := ReturnMsg{}
	customerPayment := CustomerPaymet{}

	loandetails, err := GetCustomerLoanList(custID, custUniqueID)
	if err != nil {
		returnMsg.StatusMsg = ErrNoCustomerExist
		returnMsg.Status = false
		return returnMsg
	}
	fmt.Println(loandetails)

	// TODO: login = epoint_merchant_login() from 284-306

	// check if customer balance in wallet is enough for payment
	//	logout = epoint_merchant_logout
	var tempBalance float64 //balance is temp
	if amount > tempBalance {
		returnMsg.StatusMsg = ErrNotEnoughWalletBalance
		returnMsg.Status = false
		return returnMsg
	} else {

		var currentTime = time.Now().Format("2000-03-03 15:03:01")
		stringCustID := strconv.Itoa(custID)

		refCode := GenerateRandomKey(5)
		PLreferenceCode := "PL-" + refCode
		customerPayment.FKCustomerID = &custID
		customerPayment.ReferenceCode = &PLreferenceCode
		customerPayment.PayMentAmount = &amount
		customerPayment.DatePaid = &currentTime
		customerPayment.PaidBy = &stringCustID

		// TODO: PROCESS AMOUNT IN EPOINT WALLET from 316- 333
		err = db.Model(&customerPayment).Insert()
		if err != nil {
			returnMsg.StatusMsg = ErrNotEnoughWalletBalance
			returnMsg.Status = false
			return returnMsg
		}

		// TODO: line 338-413

	}

	returnMsg.StatusMsg = MsgCustomerSuccessPaidLoan
	returnMsg.Status = true
	return returnMsg
}

// GetCustomerLoanList gets customer base on customerID and is_paid equal to 0
func GetCustomerLoanList(custID int, custUniqueID string) (CustomerLoans, error) {

	customerloans := CustomerLoans{}

	q := db.Select().
		From("tblCustomerLoan").
		Where(dbx.HashExp{"fk_customer_id": custID, "is_paid": 0}).
		OrderBy("id")
	err := q.All(&customerloans)

	if err != nil {
		message := ErrCustomerNotFound
		return customerloans, errors.New(message)
	}
	return customerloans, nil

}
