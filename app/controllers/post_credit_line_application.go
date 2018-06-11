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
	"github.com/labstack/gommon/random"
)

// LoanCreditLimit contains information of tblLoanCreditLimit
type LoanCreditLimit struct {
	ID          int
	Tier        *int
	code        *string
	Description *string
	Amount      *float64
	NoOfDays    *int
	Active      *string
	CreatedDate time.Time
}

// SystemSetting contains information of tblsystemsetting
type SystemSetting struct {
	ID           int
	Name         *string
	Code         *string
	Description  *string
	Value        *string
	SettingType  *string
	IsActive     *string
	SMSMessage   *string
	EmailMessage *string
	Subject      *string
	From         *string
	To           *string
	Cc           *string
	Bcc          *string
	UpdatedBy    *int
}

// CustomerCrediteApplication contains information of tblCustomercrediteapplication
type CustomerCrediteApplication struct {
	ID               int
	FkCustomerID     *int
	CreditLineID     *int
	CreditLineAmount *float64
	ReferenceCode    *string
	Status           *string
	ProcessedBy      *string
	ProcessedDate    *string // may need to date type
}

// CustomerInfo contains information of tblcustomerBasicinfo
/*
type CustomerInfo struct {
	ID        int
	FirstName *string `db:"first_name"`
}

// TableName gets tblcustomerbasicinfo from database
func (a CustomerInfo) TableName() string {
	return "view_customer_info"
}
*/

// TableName gets tblloancreditlimit from database
func (a LoanCreditLimit) TableName() string {
	return "tblLoancreditlimit"
}

// TableName gets tblsystemsettings from database
func (a SystemSetting) TableName() string {
	return "tblsystemsettings"
}

// TableName gets tblLoancreditlimit from database
func (a CustomerCrediteApplication) TableName() string {
	return "tblcustomercreditlineapplication"
}

var (
	// MsgSuccessApplilyCredit is given when the customer is success to add credit
	MsgSuccessApplilyCredit = "Customer successfully added a credit line."

	// ErrNoCustomerExist is given when the customer is not success to add credit
	ErrNoCustomerExist = "No existing customer"
)

// PostCreditLineApplication allows to add credit for customer in tblCustomerCreditLineApplication from database if it exists
func (co *Controllers) PostCreditLineApplication(c echo.Context) error {

	custUniqueID := c.FormValue("R1")
	// Get customer ID
	customerID := c.Get("customerID").(int)

	settingMessage := SystemSetting{}

	// Get customer information in tblcustomerbasicinfo from database if it exists
	customerinfo := CustomerInfo{}
	dbcustInfo := db.Select().
		From("view_customer_info").
		Where(dbx.HashExp{"cust_unique_id": custUniqueID})

	err := dbcustInfo.One(&customerinfo)

	if err != nil {
		message := ErrNoCustomerExist
		return SendErrorResponse(c, http.StatusBadRequest, message)
	}
	fmt.Println(*customerinfo.FirstName)

	// Get loancreditlimit
	limit, err := GetLoanCreditLimit("") //It return a struct of loancreditlimit

	if err != nil {
		message := ErrNoCustomerExist
		SendErrorResponse(c, http.StatusBadRequest, message)
	}

	fmt.Println(*limit.Active)

	// get random key for referenceCode
	referenceCode := GenerateRandomKey(5)
	fmt.Println(referenceCode)

	// Get struct of systemsetting
	creditApprovalSettings, err := GlobalSetting("1") //tblsystemsettings

	if err != nil {
		message := ErrNoCustomerExist
		SendErrorResponse(c, http.StatusBadRequest, message)
	}
	fmt.Println(creditApprovalSettings.ID)

	// Set variable for input in struct of customercrediteApplication
	tempReferenceCode := "CA-" + referenceCode
	tempProcessedBy := "SYSTEM"
	datetime := time.Now().Format("2006-01-02 15:04:05")

	customerCrediteApplication := CustomerCrediteApplication{
		FkCustomerID:     &customerID,
		CreditLineID:     &limit.ID,
		CreditLineAmount: limit.Amount,
		ReferenceCode:    &tempReferenceCode,
		ProcessedBy:      &tempProcessedBy,
		ProcessedDate:    &datetime,
	}

	//	fmt.Println(*customerCrediteApplication.ProcessedDate)

	// Set variable to customerCrediteApplication.Status
	one := "1"
	approved := "APPROVED"
	pending := "PENDING"

	settingMessage = SystemSetting{}

	// Credit Approval
	if *creditApprovalSettings.Value == one { // auto approve

		customerCrediteApplication.Status = &approved
		settingMessage, err := GlobalSetting("3") // credit approved notification, for send message in send message function

		if err != nil {
			message := ErrNoCustomerExist
			return SendErrorResponse(c, http.StatusBadRequest, message)
		}
		fmt.Println(settingMessage)

	} else {

		// Credit is not approval
		customerCrediteApplication.Status = &pending
		settingMessage, err := GlobalSetting("4") // credit pending notification, for send message in send message function

		if err != nil {
			message := ErrNoCustomerExist
			return SendErrorResponse(c, http.StatusBadRequest, message)
		}

		fmt.Println(settingMessage)
	}

	fmt.Println(*customerCrediteApplication.Status) // display approved
	fmt.Println(*customerCrediteApplication.ProcessedDate)

	// Insert value into tblCustomerCreditLineApplication from database
	err = db.Model(&customerCrediteApplication).Insert()

	if err != nil {
		message := ErrNoCustomerExist
		fmt.Println(err)
		return SendErrorResponse(c, http.StatusBadRequest, message)
	}

	// Replace amount to {{amount}} in sms message

	loanAmount := strconv.FormatFloat(*limit.Amount, 'f', 2, 64)
	p := *settingMessage.SMSMessage                              // Convert amount to string for replacement in settingMessage.SMSMessage
	SMSmessage := strings.Replace(p, "{amount}", loanAmount, -1) // -1 replace all amount
	fmt.Println(SMSmessage)

	r := strings.NewReplacer("{amount}", loanAmount, "firstname", *customerinfo.FirstName)
	r.Replace(*settingMessage.EmailMessage)
	fmt.Print(r)
	// TODO : Add sent_mail function
	// TODO : Add sent_mail function

	// It is success to apply credit
	message := MsgSuccessApplilyCredit
	return SendOKResponse(c, message)

}

// GetLoanCreditLimit return struct of loanCreditLimit for customer
func GetLoanCreditLimit(creditLineID string) (LoanCreditLimit, error) {
	loanCreditLimit := LoanCreditLimit{}
	q := db.Select().
		From("tblLoanCreditLimit")
	if len(creditLineID) > 0 {
		q = q.Where(dbx.HashExp{"id": creditLineID})

	} else {
		q = q.Where(dbx.HashExp{"active": "YES"})
	}

	err := q.One(&loanCreditLimit)

	if err != nil {
		message := ErrProblemOccured
		return loanCreditLimit, errors.New(message)
	}
	return loanCreditLimit, nil

}

// GenerateRandomKey return a random key
func GenerateRandomKey(length uint8) string {
	if length == 0 {
		length = 20
	}

	randomKey := random.New().String(length, random.Uppercase+random.Numeric)

	return randomKey
}

// GlobalSetting return a struct of systemsetting
func GlobalSetting(checkProveID string) (SystemSetting, error) {
	if len(checkProveID) == 0 {
		checkProveID = ""
	}

	systemSetting := SystemSetting{}

	ID, err := strconv.Atoi(checkProveID)
	if err != nil {
		message := ErrProblemOccured
		return systemSetting, errors.New(message)
	}

	q := db.Select().
		From("tblSystemSettings").
		Where(dbx.HashExp{"id": ID})

	err = q.One(&systemSetting)

	if err != nil {
		message := ErrProblemOccured
		return systemSetting, errors.New(message)
	}
	return systemSetting, nil

}
