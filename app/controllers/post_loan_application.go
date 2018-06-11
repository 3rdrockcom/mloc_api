package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/labstack/echo"
)

// ReturnMsg returns status information of postLoan application
type ReturnMsg struct {
	Status    bool
	StatusMsg string
}

// CustomerLoanApplication contains  information of tblcustomerloanapplication in database
type CustomerLoanApplication struct { // need to fill in more element
	ID                  int
	FkCustomerID        *int
	LoanAmount          *float64
	InterestAmount      *float64
	FeeAmount           *float64
	TotalAmount         *float64
	ReferenceCode       *string
	DueDate             *string
	LoanDate            *string
	CreatedBy           *string
	CreatedDate         *string
	Status              *string
	ProcessedBy         *string
	ProcessedDate       *string // may need to date type
	EpointTransactionID *string
}

// TableName gets tblCustomerLoanApplication from database
func (a CustomerLoanApplication) TableName() string {
	return "tblCustomerLoanApplication"
}

var (
	// ErrProcessingLoanApp is given if loan application processes error occur
	ErrProcessingLoanApp = "Error while processing loan application."

	// ErrTransferAmount  is given if loan amount cant transfer in epoint wallet
	ErrTransferAmount = "Cannot transfer amount in EPOINT. Please contact epoint admin."

	// ErrInvalidUserPassword is given if it username or password is wrong
	ErrInvalidUserPassword = "Invalid EPOINT username or password."

	// SuccessApplyLoanApplication is given if loan application process sucessful
	SuccessApplyLoanApplication = "Customer successfully applied for a loan."
)

// PostLoanApplication procedures applied for a loan
func (co *Controllers) PostLoanApplication(c echo.Context) error {
	// Get custoemrID
	customerID := c.Get("customerID").(int)
	// Get customer Unique ID and amount
	custUniqueID := c.FormValue("R1")
	amount := c.FormValue("R2")

	// Convert amount string type to float64 type.If it can't convert, then return error
	loanAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		message := ErrNoValidAmount
		return SendErrorResponse(c, http.StatusBadRequest, message)
	}

	// Call loanApplication function to procedure loan application, it returns staturs true or false.
	// If return false, display error accur,otherwise display apply loan success
	update := loanApplication(customerID, custUniqueID, loanAmount, c)

	if update.Status == false {
		message := update.StatusMsg
		return SendErrorResponse(c, http.StatusBadRequest, message)
	}

	// Loan application is success to apply
	message := update.StatusMsg
	return SendOKResponse(c, message)

}

//  loanAppicaion produces apply loan application, it return result of apply loan
func loanApplication(custID int, custUniqueID string, amount float64, c echo.Context) ReturnMsg {

	customerInfo := GetCustomer{} // GetCustomer is in post_computer_loan.go file
	returnMsg := ReturnMsg{}
	settingMessage := SystemSetting{}
	// It checks customerID is empty or not. if customerID is empty, it returns false
	q := db.Select().
		From("view_customer_info").
		Where(dbx.HashExp{"customer_id": custID}) // cust_unique_id is *string type
	err := q.One(&customerInfo)

	if err != nil {
		returnMsg.StatusMsg = ErrNoCustomerExist
		returnMsg.Status = false
		return returnMsg
	}
	//	fmt.Println(*customerInfo.FirstName)
	//	fmt.Println(customerInfo.CustomerID)
	//	fmt.Println(*customerInfo.AvailableCredit)

	if amount > *customerInfo.AvailableCredit {
		returnMsg.StatusMsg = ErrNotEnoughCredit
		returnMsg.Status = false
		return returnMsg
	}

	// Get length of 5 random number for refCode
	refCode := GenerateRandomKey(5) // in post_credit_line_application.go file
	//	fmt.Println(refCode)

	// Get loan approval message
	loanApprovalSetting, err := GlobalSetting("2") // GlobalSetting function is in post-credit_line_application
	if err != nil {
		returnMsg.StatusMsg = ErrProblemOccured
		returnMsg.Status = false
		return returnMsg
	}

	//	fmt.Println(loanApprovalSetting.ID)
	// computedLoan puts information to tblcomputerloanapplication in database
	computedLoan, err := ComputerLoanApplication(custID, custUniqueID, amount) // in post_computer_loan.go file
	fmt.Println(*computedLoan.Amount)
	fmt.Println(computedLoan.Fee)
	fmt.Println(computedLoan.TotalAmount)

	var currentTime = time.Now().Format("2000-03-03 15:03:01")
	refCode = "LA-" + refCode
	system := "SYSTEM"

	customerloanApplication := CustomerLoanApplication{
		FkCustomerID:   &custID,
		LoanAmount:     &amount,
		InterestAmount: &computedLoan.Interest,
		FeeAmount:      &computedLoan.Fee,
		TotalAmount:    &computedLoan.TotalAmount,
		ReferenceCode:  &refCode,
		DueDate:        &computedLoan.DueDate,
		LoanDate:       &currentTime,
		CreatedBy:      &system,
		CreatedDate:    &currentTime,
	}

	one := "1"
	approved := "APPROVED"
	system = "SYSTEM"

	if loanApprovalSetting.Value == &one {
		customerloanApplication.Status = &approved
		customerloanApplication.ProcessedBy = &system
		customerloanApplication.ProcessedDate = &currentTime

		settingMessage, err := GlobalSetting("6") // in Post_credit_line_application.go file
		if err != nil {
			returnMsg.StatusMsg = ErrProblemOccured
			returnMsg.Status = false
			return returnMsg
		}
		fmt.Println(settingMessage)

		//TODO: $login = epoint_merchant_login(); from 146-177
		/* forward the amount loan in the EPOINT wallet via curl starts here
		   function is inside the helper/mloc_helper.php
		*/
		//	EPOINT_MTID := "65626"
		//	EPOINT_USER := "mlocmtusr01"
		//	EPOINT_PASSWORD := "1h6W8C6H20V4"

		url := "http://epointjetdev.epointserver.com/api/merchant/merchant_login?P01={{EPOINT_MTID}}&P02={{EPOINT_USER}}&P03={{EPOINT_PASSWORD}}"

		req, _ := http.NewRequest("GET", url, nil)

		req.Header.Add("X-API-KEY", "{{customer_api_key}}")

		res, _ := http.DefaultClient.Do(req)

		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)

		fmt.Println(res)
		fmt.Println(string(body))

		// TODO:  do bellow is commenting
		/*	if req.ResponseCode == "0000" {
				// TODO: process amount in epoint wallet
				//	c.Param("session_id")
				c.Param("amount") = amount
				c.Param("reference_number") = customerloanApplication.ReferenceCode
				c.Param("source") = "p"
				c.Param("destination") = customerInfo.ProgramCustomerID
				c.Param("description") = "LOAN_APPROVED_VIA_mloc"
				c.Param("mobile") = customerInfo.MobileNumber
				transferLoan = epointMerchantFundTransfer(c) // TODO :Need to write this function

				if transferLoan.ResponseCode == "0000" {
					//TODO:
					//	customerloanApplication.EpointTransactionID :=transfer_loan['ResponseMessage']['epoint_transaction_id'];

				} else {
					returnMsg.StatusMsg = ErrTransferAmount
					returnMsg.Status = false
					return returnMsg
				}
			} else {
				returnMsg.StatusMsg = ErrInvalidUserPassword
				returnMsg.Status = false
				return returnMsg
			}
		*/
	} else {
		pending := "PENDDING"
		customerloanApplication.Status = &pending

		settingMessage, err := GlobalSetting("7")
		if err != nil {
			returnMsg.StatusMsg = ErrProblemOccured
			returnMsg.Status = false
			return returnMsg
		}
		fmt.Println(customerloanApplication)
		fmt.Println(settingMessage)
	}

	// Insert cusomerloanApplication in to tblcustomerloanapplication in database
	err = db.Model(&customerloanApplication).Insert()
	if err != nil {
		returnMsg.StatusMsg = ErrProcessingLoanApp
		returnMsg.Status = false
		return returnMsg
	}

	// Replace amount to sms message
	loanAmount := strconv.FormatFloat(amount, 'f', 2, 64)
	p := *settingMessage.SMSMessage                              // Convert amount to string for replacement in settingMessage.SMSMessage
	SMSmessage := strings.Replace(p, "{amount}", loanAmount, -1) // -1 replace all amount
	fmt.Println(SMSmessage)
	r := strings.NewReplacer("{amount}", loanAmount, "firstname", *customerInfo.FirstName)

	r.Replace(*settingMessage.EmailMessage)

	// TODO: Send_sms()
	// TODO: send_mail()

	// Loan application is success
	returnMsg.Status = true
	returnMsg.StatusMsg = SuccessApplyLoanApplication
	return returnMsg

}
