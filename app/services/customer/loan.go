package customer

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/epointpayment/mloc_api_go/app/helpers"
	"github.com/epointpayment/mloc_api_go/app/models"
	Notifications "github.com/epointpayment/mloc_api_go/app/services/notifications"
	Mail "github.com/epointpayment/mloc_api_go/app/services/notifications/mail"
	SMS "github.com/epointpayment/mloc_api_go/app/services/notifications/sms"
	"github.com/epointpayment/mloc_api_go/app/services/payments"
	"github.com/epointpayment/mloc_api_go/app/services/payments/collection"
	"github.com/epointpayment/mloc_api_go/app/services/payments/disbursement"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/labstack/gommon/log"
	"github.com/shopspring/decimal"
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

	_, err = tx.Update(
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

	// Check if credit approval is automatic or manual
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
	err = tx.Model(&customerCreditLineApplication).Insert(
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

	// Initialize notifications service
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

// ComputedLoan contains information about a loan
type ComputedLoan struct {
	AvailableCredit  decimal.Decimal `json:"available_credit"`
	Amount           decimal.Decimal `json:"amount"`
	Fee              decimal.Decimal `json:"fee"`
	Interest         decimal.Decimal `json:"interest"`
	DateApplied      string          `json:"date_applied"`
	DueDate          string          `json:"due_date"`
	DueDateFormatted string          `json:"due_date_formatted"`
	TotalAmount      decimal.Decimal `json:"total_amount"`
}

// ComputeLoanApplication calculates information about a loan
func (l *Loan) ComputeLoanApplication(baseAmount decimal.Decimal) (computed ComputedLoan, err error) {
	// Get detailed customer information
	customer, err := l.cs.Info().GetDetails()
	if err != nil {
		return
	}
	availableCredit := decimal.NewFromFloat(customer.AvailableCredit.ValueOrZero())

	// Check if requested loan amount is a valid amount
	if baseAmount.GreaterThan(availableCredit) {
		err = ErrNotEnoughAvailableCredit
		return
	}

	// Get loan fee
	fee, err := l.cs.Loan().GetFee()
	if err != nil {
		return
	}

	// Calculate fee
	feeAmount := decimal.NewFromFloat(fee.Fixed.ValueOrZero())
	if fee.Percentage.ValueOrZero() > 0.0 {
		feePercentage := decimal.NewFromFloat(fee.Percentage.ValueOrZero() / 100)
		feeAmount = baseAmount.Mul(feePercentage)
	}

	// Get loan interest
	interest, err := l.cs.Loan().GetInterest()
	if err != nil {
		return
	}

	// Calculate interest
	interestAmount := decimal.NewFromFloat(interest.Fixed.ValueOrZero())
	if interest.Percentage.Float64 > 0.0 {
		feeInterest := decimal.NewFromFloat(interest.Percentage.ValueOrZero() / 100)
		interestAmount = baseAmount.Mul(feeInterest)
	}

	// Get loan credit limit information
	creditLimit, err := l.cs.Loan().GetLoanCreditLimit(int(customer.CreditLineID.Int64))
	if err != nil {
		err = ErrLoanCreditLimitNotFound
		return
	}

	// Determine due date
	t := time.Now().UTC()
	dueDate := t.AddDate(0, 0, int(creditLimit.NumberOfDays.Int64))

	// Prepare data
	computed.AvailableCredit = availableCredit.RoundBank(helpers.DefaultCurrencyPrecision)
	computed.Amount = baseAmount.RoundBank(helpers.DefaultCurrencyPrecision)
	computed.Fee = feeAmount.RoundBank(helpers.DefaultCurrencyPrecision)
	computed.Interest = interestAmount.RoundBank(helpers.DefaultCurrencyPrecision)
	computed.DateApplied = t.Format("2006-01-02 15:04:05")
	computed.DueDate = dueDate.Format("2006-01-02 15:04:05")
	computed.DueDateFormatted = dueDate.Format("01-02-2006 03:04 PM")
	computed.TotalAmount = computed.Amount.Add(computed.Fee).Add(computed.Interest)

	return
}

// ProcessLoanApplication processes a loan application
func (l *Loan) ProcessLoanApplication(method string, bankAccountID int, baseAmount decimal.Decimal) (err error) {
	// Get detailed customer information
	customer, err := l.cs.Info().GetDetails()
	if err != nil {
		return
	}
	availableCredit := decimal.NewFromFloat(customer.AvailableCredit.ValueOrZero())

	// Check if requested loan amount is a valid amount
	if baseAmount.GreaterThan(availableCredit) {
		err = ErrNotEnoughAvailableCredit
		return
	}

	// Check if loan approval is automatic or manual
	loanApproval, err := l.cs.Settings().Get(2)
	if err != nil {
		return
	}

	// Generate reference code
	refCode := loanApproval.Code.String + "-" + generateRandomKey(5)

	// Calculate loan application
	computed, err := l.ComputeLoanApplication(baseAmount)
	if err != nil {
		return
	}

	// Determine loan date
	t := time.Now().UTC()

	loanAmount, _ := baseAmount.Float64()
	loanInterest, _ := computed.Interest.Float64()
	loanFee, _ := computed.Fee.Float64()
	loanTotal, _ := computed.TotalAmount.Float64()

	customerLoanApplication := models.CustomerLoanApplication{
		CustomerID:     null.IntFrom(int64(customer.ID)),
		LoanAmount:     null.FloatFrom(loanAmount),
		InterestAmount: null.FloatFrom(loanInterest),
		FeeAmount:      null.FloatFrom(loanFee),
		TotalAmount:    null.FloatFrom(loanTotal),
		ReferenceCode:  null.StringFrom(refCode),
		DueDate:        null.StringFrom(computed.DueDate),
		LoanDate:       null.StringFrom(t.Format("2006-01-02 15:04:05")),
		CreatedBy:      null.StringFrom("SYSTEM"),
		CreatedDate:    null.StringFrom(t.Format("2006-01-02 15:04:05")),
	}

	// Determine application status and notification
	customerLoanApplication.Status = null.StringFrom("PENDING")
	systemSettingID := 7

	if loanApproval.Value.String == "1" {
		customerLoanApplication.Status = null.StringFrom("APPROVED")
		customerLoanApplication.ProcessedBy = null.StringFrom("SYSTEM")
		customerLoanApplication.ProcessedDate = null.StringFrom(t.Format("2006-01-02 15:04:05"))
		systemSettingID = 6

		// Get bank account
		customerBankAccount := new(models.CustomerBankAccount)

		if bankAccountID > 0 {
			customerBankAccount, err = l.cs.BankAccount().Get(bankAccountID)
			if err != nil {
				return
			}
		}

		// Initialize payment service
		ps := payments.New()
		if err != nil {
			return
		}

		// Prepare disbursement request
		disbursementRequest := disbursement.Request{
			Method:                  method,
			Customer:                *customer,
			CustomerBankAccount:     *customerBankAccount,
			CustomerLoanApplication: customerLoanApplication,
			Description:             "Loan_approved_via_MLOC",
		}
		disbursementResponse := disbursement.Response{}

		// Execute payment disbursement
		disbursementResponse, err = ps.Disbursement(disbursementRequest)
		if err != nil {
			err = ErrIssuerFailedTransfer
			return
		}

		// Set transaction information
		customerLoanApplication.EpointTransactionID = null.StringFrom(disbursementResponse.TransactionID)
	}

	tx, err := DB.Begin()
	if err != nil {
		return
	}

	// Store results
	err = tx.Model(&customerLoanApplication).Insert()
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

	// Initialize notifications service
	sn := Notifications.New()
	if err != nil {
		return
	}

	// Prepare template token replacer
	r := strings.NewReplacer(
		"{amount}", baseAmount.StringFixed(helpers.DefaultCurrencyPrecision),
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

	return
}

// ProcessLoanPayment processes a loan payment
func (l *Loan) ProcessLoanPayment(paymentAmount decimal.Decimal) (err error) {
	// Get detailed customer information
	customer, err := l.cs.Info().GetDetails()
	if err != nil {
		return
	}
	// Generate reference code
	refCode := "PL" + "-" + generateRandomKey(5)

	// Determine payment date
	t := time.Now().UTC()

	// Prepare loan payment information
	pa, _ := paymentAmount.Float64()
	customerPayment := models.CustomerPayment{
		CustomerID:    null.IntFrom(int64(customer.ID)),
		ReferenceCode: null.StringFrom(refCode),
		PaymentAmount: null.FloatFrom(pa),
		DatePaid:      null.StringFrom(t.Format("2006-01-02 15:04:05")),
		PaidBy:        null.StringFrom(strconv.FormatInt(int64(customer.ID), 10)),
	}

	// Initialize payments service
	ps := payments.New()
	if err != nil {
		return
	}

	// Prepare collection request
	collectionRequest := collection.Request{
		Method:          payments.MethodEPOINT,
		Customer:        *customer,
		CustomerPayment: customerPayment,
		Description:     "Loan_approved_via_MLOC",
	}
	collectionResponse := collection.Response{}

	// Execute payment collection
	collectionResponse, err = ps.Collection(collectionRequest)
	if err != nil {
		err = ErrIssuerFailedTransfer
		return
	}

	// Set transaction information
	customerPayment.EpointTransactionID = null.StringFrom(collectionResponse.TransactionID)

	tx, err := DB.Begin()
	if err != nil {
		return
	}

	// Store loan payment information
	err = tx.Model(&customerPayment).Insert()
	if err != nil {
		tx.Rollback()
		return
	}

	// Get list of unpaid (active) customer loans
	loanList, err := l.GetCustomerLoanList()
	if err != nil {
		return
	}

	// Set payment amount left to distribute among loans
	paymentAmountBalance := paymentAmount
	for _, loanEntry := range loanList {
		if paymentAmountBalance.LessThanOrEqual(decimal.Zero) {
			continue
		}
		isPaid := 0
		loanTotalPaidAmount := decimal.NewFromFloat(loanEntry.TotalPaidAmount.ValueOrZero())

		loanPrincipalAmount := decimal.NewFromFloat(loanEntry.LoanAmount.ValueOrZero())
		loanPrincipalAmountPaid := decimal.NewFromFloat(loanEntry.TotalPaidPrincipal.ValueOrZero())
		loanPrincipalAmountUnpaid := loanPrincipalAmount.Sub(loanPrincipalAmountPaid)
		paymentPrincipalApplied := decimal.Zero

		loanFeeAmount := decimal.NewFromFloat(loanEntry.FeeAmount.ValueOrZero())
		loanFeeAmountPaid := decimal.NewFromFloat(loanEntry.TotalPaidFee.ValueOrZero())
		loanFeeAmountUnpaid := loanFeeAmount.Sub(loanFeeAmountPaid)
		paymentFeeApplied := decimal.Zero

		// Use available payment balance to pay loan fee
		if loanFeeAmountUnpaid.GreaterThan(decimal.Zero) {
			// Pay loan fee with payment balance
			if paymentAmountBalance.GreaterThanOrEqual(loanFeeAmountUnpaid) {
				// Pay loan fee entirely
				paymentFeeApplied = loanFeeAmountUnpaid
				paymentAmountBalance = paymentAmountBalance.Sub(paymentFeeApplied)
			} else {
				// Pay loan fee with remaining payment balance
				paymentFeeApplied = paymentAmountBalance
				paymentAmountBalance = decimal.Zero
			}
		}

		// Use available payment balance to pay loan principal
		if loanPrincipalAmountUnpaid.GreaterThan(decimal.Zero) {
			// Pay loan principal with payment balance
			if paymentAmountBalance.GreaterThanOrEqual(loanPrincipalAmountUnpaid) {
				// Pay loan principal entirely
				paymentPrincipalApplied = loanPrincipalAmountUnpaid
				paymentAmountBalance = paymentAmountBalance.Sub(paymentPrincipalApplied)
			} else {
				// Pay loan principal with remaining payment balance
				paymentPrincipalApplied = paymentAmountBalance
				paymentAmountBalance = decimal.Zero
			}
		}

		// Check if loan has been paid off
		totalPaidFee := loanFeeAmountPaid.Add(paymentFeeApplied)
		totalPaidPrincipal := loanPrincipalAmountPaid.Add(paymentPrincipalApplied)
		if totalPaidPrincipal.Equal(loanPrincipalAmount) && totalPaidFee.Equal(loanFeeAmount) {
			isPaid = 1
		}

		// Convert to float
		totalPaidFeeAmount, _ := totalPaidFee.Float64()
		totalPaidPrincipalAmount, _ := totalPaidPrincipal.Float64()
		totalPaidAmount, _ := loanTotalPaidAmount.Add(paymentPrincipalApplied).Add(paymentFeeApplied).Float64()

		// Update customer loan information
		customerLoan := models.CustomerLoan{
			ID:                 loanEntry.ID,
			TotalPaidPrincipal: null.FloatFrom(totalPaidPrincipalAmount),
			TotalPaidFee:       null.FloatFrom(totalPaidFeeAmount),
			IsPaid:             null.IntFrom(int64(isPaid)),
			TotalPaidAmount:    null.FloatFrom(totalPaidAmount),
		}
		err = tx.Model(&customerLoan).Update(
			"TotalPaidPrincipal",
			"TotalPaidFee",
			"IsPaid",
			"TotalPaidAmount",
		)
		if err != nil {
			tx.Rollback()
			return
		}

		// Convert to float
		settlementAmount, _ := paymentPrincipalApplied.Add(paymentFeeApplied).Float64()
		principalAmount, _ := paymentPrincipalApplied.Float64()
		feeAmount, _ := paymentFeeApplied.Float64()

		// Insert settlement information
		customerSettlement := models.CustomerSettlement{
			CustomerID:        null.IntFrom(int64(customer.ID)),
			CustomerLoanID:    null.IntFrom(int64(loanEntry.ID)),
			CustomerPaymentID: null.IntFrom(int64(customerPayment.ID)),
			SettlementAmount:  null.FloatFrom(settlementAmount),
			PrincipalAmount:   null.FloatFrom(principalAmount),
			FeeAmount:         null.FloatFrom(feeAmount),
			CreatedDate:       null.StringFrom(t.Format("2006-01-02 15:04:05")),
		}
		err = tx.Model(&customerSettlement).Insert()
		if err != nil {
			tx.Rollback()
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	//

	// Get appropriate notification details
	systemSettingID := 9
	notification, err := l.cs.Settings().Get(systemSettingID)
	if err != nil {
		return
	}

	// Initialize notifications service
	sn := Notifications.New()
	if err != nil {
		return
	}

	// Prepare template token replacer
	r := strings.NewReplacer(
		"{amount}", paymentAmount.StringFixed(helpers.DefaultCurrencyPrecision),
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

// GetFee gets loan fee information
func (l *Loan) GetFee() (fee models.Fee, err error) {
	err = DB.Select().
		From(fee.TableName()).
		Where(dbx.HashExp{"active": "YES"}).
		One(&fee)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrLoanFeeNotFound
		}
		return
	}

	return
}

// GetInterest gets loan interest information
func (l *Loan) GetInterest() (interest models.Interest, err error) {
	err = DB.Select().
		From(interest.TableName()).
		Where(dbx.HashExp{"active": "YES"}).
		One(&interest)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrLoanInterestNotFound
		}
		return
	}

	return
}

// GetCustomerLoanList gets list of unpaid (active) customer loans
func (l *Loan) GetCustomerLoanList() (customerLoans models.CustomerLoans, err error) {
	err = DB.Select().
		From(models.CustomerLoan{}.TableName()).
		Where(dbx.HashExp{
			"fk_customer_id": l.cs.CustomerID,
			"is_paid":        0,
		}).
		OrderBy("id").
		All(&customerLoans)
	if err != nil {
		return
	}

	return
}
