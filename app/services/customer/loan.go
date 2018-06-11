package customer

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/gommon/log"

	"github.com/epointpayment/mloc_api_go/app/models"
	Notifications "github.com/epointpayment/mloc_api_go/app/services/notifications"
	Mail "github.com/epointpayment/mloc_api_go/app/services/notifications/mail"
	SMS "github.com/epointpayment/mloc_api_go/app/services/notifications/sms"

	dbx "github.com/go-ozzo/ozzo-dbx"
	null "gopkg.in/guregu/null.v3"
)

// Loan manages customer loan information
type Loan struct {
	cs *CustomerService
}

// GetTransactionHistoryByType gets all loan transactions by a particular type for a customer
func (l *Loan) GetTransactionHistoryByType(transactionType string) (transactionHistory models.TransactionsHistory, err error) {
	err = DB.Select().
		From(transactionHistory.TableName()).
		Where(dbx.HashExp{
			"fk_customer_id": l.cs.CustomerID,
			"t_type":         transactionType,
		}).
		OrderBy("t_date DESC").
		All(&transactionHistory)
	if err != nil {
		return
	}
	if len(transactionHistory) == 0 {
		err = ErrTransactionNotFound
		return
	}

	return
}

// GetTransactionHistory gets all loan transactions for a customer
func (l *Loan) GetTransactionHistory() (transactionHistory models.TransactionsHistory, err error) {
	err = DB.Select().
		From(transactionHistory.TableName()).
		Where(dbx.HashExp{"fk_customer_id": l.cs.CustomerID}).
		OrderBy("t_date DESC").
		All(&transactionHistory)
	if err != nil {
		return
	}
	if len(transactionHistory) == 0 {
		err = ErrTransactionNotFound
		return
	}

	return
}

// GetCustomerLoanTotal gets all loan transactions for a customer
func (l *Loan) GetCustomerLoanTotal() (loanTotal models.LoanTotal, err error) {
	err = DB.Select().
		From(loanTotal.TableName()).
		Where(dbx.HashExp{"fk_customer_id": l.cs.CustomerID}).
		One(&loanTotal)

	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrLoanNotFound
		}
		return
	}

	return
}

// AcceptedCustomerAgreement sets customer agreement terms and conditions flag to 1
func (l *Loan) AcceptedCustomerAgreement() (err error) {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	_, err = DB.Update(
		models.CustomerAgreement{}.TableName(),
		dbx.Params{"term_and_condition": 1},
		dbx.HashExp{"fk_customer_id": l.cs.CustomerID},
	).Execute()
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// ProcessCreditLineApplication processes the credit line application
func (l *Loan) ProcessCreditLineApplication() (isApproved bool, err error) {
	// Get list of active loan credit limits
	limit, err := l.GetLoanCreditLimit(0)
	if err != nil {
		return
	}

	// Check if credit is auto approve or manual
	creditApproval, err := l.cs.Settings().Get(1)
	if err != nil {
		return
	}

	// Generate reference code
	refCode := creditApproval.Code.String + "-" + generateRandomKey(5)

	// Determine application status and notification
	status := "PENDING"
	systemSettingID := 4

	if creditApproval.Value.String == "1" {
		status = "APPROVED"
		systemSettingID = 3
	}

	// Prepare application results
	customerCreditLineApplication := models.CustomerCreditLineApplication{
		CustomerID:       null.IntFrom(int64(l.cs.CustomerID)),
		CreditLineID:     null.IntFrom(int64(limit.ID)),
		CreditLineAmount: limit.Amount,
		Status:           null.StringFrom(status),
		ReferenceCode:    null.StringFrom(refCode),
		ProcessedBy:      null.StringFrom("SYSTEM"),
		ProcessedDate:    null.StringFrom(time.Now().UTC().Format("2006-01-02 15:04:05")),
	}

	tx, err := DB.Begin()
	if err != nil {
		return
	}

	// Store results
	err = DB.Model(&customerCreditLineApplication).Insert(
		"CustomerID",
		"CreditLineID",
		"CreditLineAmount",
		"ReferenceCode",
		"Status",
		"ProcessedBy",
		"ProcessedDate",
	)
	if err != nil {
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	// Get appropriate notification details
	notification, err := l.cs.Settings().Get(systemSettingID)
	if err != nil {
		return
	}

	// Get detailed customer information
	customer, err := l.cs.Info().GetDetails()
	if err != nil {
		return
	}

	// Initialize customer service
	sn := Notifications.New()
	if err != nil {
		return
	}

	// Prepare template token replacer
	r := strings.NewReplacer(
		"{amount}", strconv.FormatFloat(limit.Amount.Float64, 'f', -1, 64),
		"{firstname}", customer.FirstName.String,
	)

	// Prepare sms notification
	smsPayload := SMS.New()
	smsPayload.To = customer.MobileNumber.String
	smsPayload.Body = r.Replace(notification.SMSMessage.String)

	// Send sms notification
	go func() {
		err := sn.Send(smsPayload)
		if err != nil {
			log.Error(err)
		}
	}()

	// Prepare email notification
	emailPayload := Mail.New()
	emailPayload.To = append(emailPayload.To, Mail.Address{
		Address: customer.Email.String,
	})
	emailPayload.Subject = notification.Subject.String
	emailPayload.Body = r.Replace(notification.EmailMessage.String)

	// Send email notification
	go func() {
		err := sn.Send(emailPayload)
		if err != nil {
			log.Error(err)
		}
	}()

	isApproved = true
	return
}

// GetLoanCreditLimit gets the  loan credit limit for customer
func (l *Loan) GetLoanCreditLimit(id int) (loanCreditLimit models.LoanCreditLimit, err error) {
	q := DB.Select().
		From(loanCreditLimit.TableName())

	if id > 0 {
		q = q.Where(dbx.HashExp{"id": id})
	} else {
		q = q.Where(dbx.HashExp{"active": "YES"}).
			OrderBy("id ASC")
	}

	err = q.One(&loanCreditLimit)
	if err != nil {
		return
	}

	return
}
