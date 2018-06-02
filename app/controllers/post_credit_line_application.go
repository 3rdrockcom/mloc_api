package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
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

// CustomerBasic contains information of tblcustomerBasicinfo
type CustomerBasic struct {
	ID        int
	FirstName *string `db:"first_name"`
}

// TableName gets tblcustomerbasicinfo from database
func (a CustomerBasic) TableName() string {
	return "tblcustomerbasicinfo"
}

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
	return "tblcustomercreditLineapplication"
}

var (
	// MsgSuccessApplilyCredit is given when the customer is success to add credit
	MsgSuccessApplilyCredit = "Customer successfully added a credit line."

	// ErrNoCustomerExist is given when the customer is not success to add credit
	ErrNoCustomerExist = "No existing customer"
)

// PostCreditLineApplication allows to add credit in tblCustomerCreditLineApplication from database if it exists
func (co *Controllers) PostCreditLineApplication(c echo.Context) error {

	// Get customer ID
	customerID := c.Get("customerID").(int)

	// Get first name of customer base on customerID in tblcustomerbasicinfo from database if it exists
	customerinfo := CustomerBasic{}
	dbcustInfo := db.Select("first_name").
		From("tblcustomerbasicinfo").
		Where(dbx.HashExp{"id": customerID})

	err := dbcustInfo.One(&customerinfo)

	if err != nil {
		message := ErrProblemOccured
		return SendErrorResponse(c, http.StatusBadRequest, message)
	}
	fmt.Println(*customerinfo.FirstName)

	// Get struct of loancreditlimit
	limit, err := GetLoanCreditLimit("") //It return a struct of loancreditlimit

	if err != nil {
		message := ErrNoCustomerExist
		SendErrorResponse(c, http.StatusBadRequest, message)
	}

	fmt.Println(*limit.Active)
	// get random key
	referenceCode := generateRandomKey(5)
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

	fmt.Println(*customerCrediteApplication.ProcessedDate)

	// Set variable to customerCrediteApplication.Status
	one := "1"
	approved := "APPROVED"
	pending := "PENDING"
	/*****************************TODO*******************
		//settingMessage := SystemSetting{} // TODO: need to do it ,  for send message in send message function
	     *****************************************************/
	if *creditApprovalSettings.Value == one { // auto approve
		customerCrediteApplication.Status = &approved
		/*************************TODo*********************
		//	settingMessage, err := GlobalSetting("3") //TODO: credit approved notification, for send message in send message function
		**************************************************/
		if err != nil {
			message := ErrNoCustomerExist
			return SendErrorResponse(c, http.StatusBadRequest, message)
		}
		//	fmt.Println(settingMessage)

	} else {
		customerCrediteApplication.Status = &pending
		/***********************************TODO*********************
		//	settingMessage, err := GlobalSetting("4") //TODO: credit pending notification, for send message in send message function
		*****************************************************/
		if err != nil {
			message := ErrNoCustomerExist
			return SendErrorResponse(c, http.StatusBadRequest, message)
		}
		//	fmt.Println(settingMessage)
	}

	fmt.Println(*customerCrediteApplication.Status) // display approved
	fmt.Println(*customerCrediteApplication.ProcessedDate)

	// insert data into tblCustomerCreditLineApplication from database
	err = db.Model(&customerCrediteApplication).Insert()
	if err != nil {
		message := ErrNoCustomerExist
		return SendErrorResponse(c, http.StatusBadRequest, message)
	}
	/******************** TODO convertion and continue to finish it****************

	amount := strconv.FormatFloat(*limit.Amount, 'f', 2, 64)
	SMSmessage := strings.Replace(*settingMessage.SMSMessage, "{amount}", amount, -1)
	fmt.Println(SMSmessage)


			r := strings.NewReplacer("{amount}", amount,
				"firstname", *customerinfo.FirstName)

			r.Replace(*settingMessage.EmailMessage)

			// TODO : Add sent_mail function
			// TODO : Add sent_mail function
			message := MsgSuccessApplilyCredit

		return SendOKResponse(c, message)
	*/
	return nil
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

// generateRandomKey return a random key
func generateRandomKey(length uint8) string {
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
