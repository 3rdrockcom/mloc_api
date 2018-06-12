package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/labstack/echo"
)

// CustomerLoans array contains CustomerLoanstruct
type CustomerLoans []CustomerLoan

// CustomerLoan contains information about fkCustomerID
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

type CustomerSettlement struct {
	ID                int
	FKCustomerID      *int
	CustomerLoanID    *int
	CustomerPaymentID *int
	SettlementAmount  *float64
	PrincipalAmount   *float64
	FeeAmount         *float64
	CreatedDate       *string // may need to change time type
}

// TableName gets tblcustoemrpayment from database
func (a CustomerPaymet) TableName() string {
	return "tblcustomerpayment"
}

// TableName gets tblcustomerloan from database
func (a CustomerLoan) TableName() string {
	return "tblcustomerloan"
}

// TableName gets tblcustomerloan from database
func (a CustomerSettlement) TableName() string {
	return "tblcustomersettlement"
}

var (
	// ErrNotEnoughWalletBalance is given if the customer doesn't have enought wallet balance
	ErrNotEnoughWalletBalance = "You don't have enough available wallet balance."

	// ErrProcessLoanPayment is given while process pay loan
	ErrProcessLoanPayment = "Error while processing loan payment."

	// MsgCustomerSuccessPaidLoan is given while customer paid loan is successsfully
	MsgCustomerSuccessPaidLoan = "Customer successfully paid a loan."
)

// postPayLoan displays result of payloan of customer
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

// payLoanProcess processes the customer pay loan and returns status message
func payLoanProcess(custID int, custUniqueID string, amount float64) ReturnMsg {

	returnMsg := ReturnMsg{}            // Initial empty returnMsg
	customerPayment := CustomerPaymet{} // Initial empty customerPayment struct
	customerInfo := CustomerInfo{}      // Initial empty customerInfo struct
	loanUpdate := CustomerLoan{}        // Initial empty CustomerLoan struct for update tblCustomerLoan table in database
	settlement := CustomerSettlement{}  //Initial empty customersettlement for update tblcustomersettlement table in database
	// Get customer informations in view_customer_info table from database
	customerInfo, err := GetCustomerInDetail(custUniqueID)
	if err != nil {
		returnMsg.StatusMsg = ErrNoCustomerExist
		returnMsg.Status = false
		return returnMsg
	}

	// Get customer loan information in tblcustomerloan table from database
	loanDetails, err := GetCustomerLoanList(custID, custUniqueID)
	if err != nil {
		returnMsg.StatusMsg = ErrNoCustomerExist
		returnMsg.Status = false
		return returnMsg
	}

	// TODO: login = epoint_merchant_login() from 284-306

	// check if customer balance in wallet is enough for payment

	var tempBalance float64 //balance is temp
	if amount > tempBalance {
		returnMsg.StatusMsg = ErrNotEnoughWalletBalance
		returnMsg.Status = false
		return returnMsg
	}

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
		returnMsg.StatusMsg = ErrProcessLoanPayment
		returnMsg.Status = false
		return returnMsg
	}

	// Convert amount to string for replacement in settingMessage.SMSMessage
	loanAmount := strconv.FormatFloat(amount, 'f', 2, 64)

	//TODO: customerPaymentID := insert_id()
	settingMessage, err := GlobalSetting("9")
	if err != nil {
		returnMsg.StatusMsg = ErrProcessLoanPayment
		returnMsg.Status = false
		return returnMsg
	}
	p := *settingMessage.SMSMessage

	SMSmessage := strings.Replace(p, "{amount}", loanAmount, -1) // -1 replace all amount
	fmt.Println(SMSmessage)
	r := strings.NewReplacer("{amount}", loanAmount, "firstname", *customerInfo.FirstName)
	r.Replace(*settingMessage.EmailMessage)
	// TODO send_sms()
	// TODO: Send_mail()

	paymentAmount := amount
	fmt.Println(paymentAmount)

	for _, row := range loanDetails { //TODO:  foreach ($loan_details as $row)foreach

		if paymentAmount > 0 {
			fee := 0.00
			baseAmount := 0.00
			ispaind := 0
			totalFee := 0
			totalBase := 0

			fmt.Println(ispaind)
			fmt.Println(totalFee)
			fmt.Println(totalBase)

			// Check if fee is already paid in tblCustomerLoan
			if *row.FeeAmount != *row.TotalPaidFee {
				fee := *row.FeeAmount - *row.TotalPaidFee

				if paymentAmount >= fee {
					paymentAmount = paymentAmount - fee
				} else {
					fee = paymentAmount
					paymentAmount = 0
				}
			}
			isPaid := 1
			// NOTE: after tblCustomerLoan update,the data in tblCustomerCreditLine and tblCustomerLoanTotal will
			// automatically update using the mysql trigger trg_after_tblCustomerLoan_update

			// Update tblCustomerLoan
			*loanUpdate.TotalPaidPrincipal = *row.TotalPaidPrincipal + baseAmount
			*loanUpdate.TotalPaidFee = *row.TotalPaidFee + fee
			*loanUpdate.IsPaid = isPaid
			*loanUpdate.TotalPaidAmount = *row.TotalPaidAmount + (fee + baseAmount)
			err = db.Model(&loanUpdate).Update() //TODO: specific update id
			if err != nil {
				returnMsg.StatusMsg = ErrProcessLoanPayment
				returnMsg.Status = false
				return returnMsg
			}

			date := time.Now().Format("2006-01-02 15:04:05")

			// Update settlement
			*settlement.FKCustomerID = customerInfo.CustomerID
			*settlement.CustomerLoanID = row.ID
			//TODO: 	settlement.CustomerPaymentID = customer
			*settlement.SettlementAmount = fee + baseAmount
			*settlement.PrincipalAmount = baseAmount
			*settlement.FeeAmount = fee
			settlement.CreatedDate = &date

			err = db.Model(&settlement).Update()
			if err != nil {
				returnMsg.StatusMsg = ErrProcessLoanPayment
				returnMsg.Status = false
				return returnMsg
			}

		}
	}

	returnMsg.StatusMsg = MsgCustomerSuccessPaidLoan
	returnMsg.Status = true
	return returnMsg

}

// GetCustomerLoanList gets customer base on customerID and is_paid equal to 0
func GetCustomerLoanList(custID int, custUniqueID string) (CustomerLoans, error) {

	customerloans := CustomerLoans{}

	q := db.Select().
		From("tblcustomerloan").
		Where(dbx.HashExp{"fk_customer_id": custID, "is_paid": 0}).
		OrderBy("id")
	err := q.All(&customerloans)

	if err != nil {
		message := ErrCustomerNotFound
		return customerloans, errors.New(message)
	}
	return customerloans, nil

}

// GetCustomerInDetail gets information from view_customer_info table from database. If it exists, return value. Otherwise, it returns error
func GetCustomerInDetail(custUniqueID string) (CustomerInfo, error) {
	customerInfo := CustomerInfo{}

	err := db.Select().
		From("view_customer_info").
		Where(dbx.HashExp{"cust_unique_id": custUniqueID}).
		One(&customerInfo)

	if err != nil {
		message := ErrCustomerNotFound
		return customerInfo, errors.New(message)
	}

	return customerInfo, nil
}
