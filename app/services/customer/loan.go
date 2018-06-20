package customer

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/gommon/log"

	"github.com/epointpayment/mloc_api_go/app/models"
	EPOINT "github.com/epointpayment/mloc_api_go/app/services/epoint"
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
	AvailableCredit  float64 `json:"available_credit"`
	Amount           float64 `json:"amount"`
	Fee              float64 `json:"fee"`
	Interest         float64 `json:"interest"`
	DateApplied      string  `json:"date_applied"`
	DueDate          string  `json:"due_date"`
	DueDateFormatted string  `json:"due_date_formatted"`
	TotalAmount      float64 `json:"total_amount"`
}

// ComputeLoanApplication calculates information about a loan
func (l *Loan) ComputeLoanApplication(baseAmount float64) (computed ComputedLoan, err error) {
	// Get detailed customer information
	customer, err := l.cs.Info().GetDetails()
	if err != nil {
		return
	}
	availableCredit := customer.AvailableCredit.ValueOrZero()

	// Check if requested loan amount is a valid amount
	if baseAmount > availableCredit {
		err = ErrNotEnoughAvailableCredit
		return
	}

	// Get loan fee
	fee, err := l.cs.Loan().GetFee()
	if err != nil {
		return
	}

	// Calculate fee
	feeAmount := fee.Fixed.ValueOrZero()
	if fee.Percentage.ValueOrZero() > 0.0 {
		feeAmount = (fee.Percentage.ValueOrZero() / 100) * baseAmount
	}

	// Get loan interest
	interest, err := l.cs.Loan().GetInterest()
	if err != nil {
		return
	}

	// Calculate interest
	interestAmount := interest.Fixed.ValueOrZero()
	if interest.Percentage.Float64 > 0.0 {
		interestAmount = (interest.Percentage.ValueOrZero() / 100) * baseAmount
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
	computed.AvailableCredit = availableCredit
	computed.Amount = numberFormat(baseAmount, 2)
	computed.Fee = numberFormat(feeAmount, 2)
	computed.Interest = numberFormat(interestAmount, 2)
	computed.DateApplied = t.Format("2006-01-02 15:04:05")
	computed.DueDate = dueDate.Format("2006-01-02 15:04:05")
	computed.DueDateFormatted = dueDate.Format("01-02-2006 03:04 PM")
	computed.TotalAmount = computed.Amount + computed.Fee + computed.Interest

	return
}

// ProcessLoanApplication processes a loan application
func (l *Loan) ProcessLoanApplication(baseAmount float64) (err error) {
	// Get detailed customer information
	customer, err := l.cs.Info().GetDetails()
	if err != nil {
		return
	}
	availableCredit := customer.AvailableCredit.ValueOrZero()

	// Check if requested loan amount is a valid amount
	if baseAmount > availableCredit {
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

	customerLoanApplication := models.CustomerLoanApplication{
		CustomerID:     null.IntFrom(int64(customer.ID)),
		LoanAmount:     null.FloatFrom(baseAmount),
		InterestAmount: null.FloatFrom(computed.Interest),
		FeeAmount:      null.FloatFrom(computed.Fee),
		TotalAmount:    null.FloatFrom(computed.TotalAmount),
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

		// Initialze epoint service
		es := new(EPOINT.EpointService)
		es, err = EPOINT.New()
		if err != nil {
			return
		}

		// Login to epoint service
		_, err = es.GetLogin()
		if err != nil {
			err = ErrIssuerInvalidUserPassword
			return
		}

		// Transfer funds from prefund to customer wallet using epoint service
		fundTransferRequest := EPOINT.FundTransferRequest{
			Amount:          baseAmount,
			ClientReference: customerLoanApplication.ReferenceCode.String,
			Source:          "P",
			Destination:     strconv.FormatInt(customer.ProgramCustomerID.Int64, 10),
			Description:     "Loan_approved_via_MLOC",
			MobileNumber:    customer.MobileNumber.String,
		}
		ft := EPOINT.FundTransferResponse{}
		ft, err = es.GetFundTransfer(fundTransferRequest)
		if err != nil {
			err = ErrIssuerFailedTransfer
			return
		}

		customerLoanApplication.EpointTransactionID = null.StringFrom(ft.TransactionID)
	}

	tx, err := DB.Begin()
	if err != nil {
		return
	}

	// Store results
	err = DB.Model(&customerLoanApplication).Insert()
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
		"{amount}", strconv.FormatFloat(baseAmount, 'f', -1, 64),
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
func (l *Loan) ProcessLoanPayment(paymentAmount float64) (err error) {
	// Get detailed customer information
	customer, err := l.cs.Info().GetDetails()
	if err != nil {
		return
	}

	// Initialze epoint service
	es := new(EPOINT.EpointService)
	es, err = EPOINT.New()
	if err != nil {
		return
	}

	// Login to epoint service
	_, err = es.GetLogin()
	if err != nil {
		err = ErrIssuerInvalidUserPassword
		return
	}

	// Get customer balance information
	customerBalanceRequest := EPOINT.CustomerBalanceRequest{
		CustomerID:   int(customer.ProgramCustomerID.Int64),
		MobileNumber: customer.MobileNumber.String,
	}
	cb, err := es.GetCustomerBalance(customerBalanceRequest)
	if err != nil {
		err = ErrIssuerUnableToAccessBalance
		return
	}

	// Check if there is enough funds available in wallet for payment
	if paymentAmount > cb.AvailableBalance {
		err = ErrIssuerInsufficientFunds
		return
	}

	// Generate reference code
	refCode := "PL" + "-" + generateRandomKey(5)

	// Determine payment date
	t := time.Now().UTC()

	// Transfer funds from customer wallet to settlement using epoint service
	fundTransferRequest := EPOINT.FundTransferRequest{
		Amount:          paymentAmount,
		ClientReference: refCode,
		Source:          strconv.FormatInt(customer.ProgramCustomerID.Int64, 10),
		Destination:     "S",
		Description:     "Loan_payment_via_MLOC",
		MobileNumber:    customer.MobileNumber.String,
	}
	ft := EPOINT.FundTransferResponse{}
	ft, err = es.GetFundTransfer(fundTransferRequest)
	if err != nil {
		err = ErrIssuerFailedTransfer
		return
	}

	tx, err := DB.Begin()
	if err != nil {
		return
	}

	// Prepare loan payment information
	customerPayment := models.CustomerPayment{
		CustomerID:          null.IntFrom(int64(customer.ID)),
		ReferenceCode:       null.StringFrom(refCode),
		PaymentAmount:       null.FloatFrom(paymentAmount),
		DatePaid:            null.StringFrom(t.Format("2006-01-02 15:04:05")),
		PaidBy:              null.StringFrom(strconv.FormatInt(int64(customer.ID), 10)),
		EpointTransactionID: null.StringFrom(ft.TransactionID),
	}
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
		if paymentAmountBalance <= 0.0 {
			continue
		}
		isPaid := 0
		loanTotalPaidAmount := loanEntry.TotalPaidAmount.ValueOrZero()

		loanPrincipalAmount := loanEntry.LoanAmount.ValueOrZero()
		loanPrincipalAmountPaid := loanEntry.TotalPaidPrincipal.ValueOrZero()
		loanPrincipalAmountUnpaid := loanPrincipalAmount - loanPrincipalAmountPaid
		paymentPrincipalApplied := 0.0

		loanFeeAmount := loanEntry.FeeAmount.ValueOrZero()
		loanFeeAmountPaid := loanEntry.TotalPaidFee.ValueOrZero()
		loanFeeAmountUnpaid := loanFeeAmount - loanFeeAmountPaid
		paymentFeeApplied := 0.0

		// Use available payment balance to pay loan fee
		if loanFeeAmountUnpaid > 0.0 {
			// Pay loan fee with payment balance
			if paymentAmountBalance >= loanFeeAmountUnpaid {
				// Pay loan fee entirely
				paymentFeeApplied = loanFeeAmountUnpaid
				paymentAmountBalance -= paymentFeeApplied
			} else {
				// Pay loan fee with remaining payment balance
				paymentFeeApplied = paymentAmountBalance
				paymentAmountBalance = 0
			}
		}

		// Use available payment balance to pay loan principal
		if loanPrincipalAmountUnpaid > 0.0 {
			// Pay loan principal with payment balance
			if paymentAmountBalance >= loanPrincipalAmountUnpaid {
				// Pay loan principal entirely
				paymentPrincipalApplied = loanPrincipalAmountUnpaid
				paymentAmountBalance -= paymentPrincipalApplied
			} else {
				// Pay loan principal with remaining payment balance
				paymentPrincipalApplied = paymentAmountBalance
				paymentAmountBalance = 0
			}
		}

		// Check if loan has been paid off
		totalPaidFee := loanFeeAmountPaid + paymentFeeApplied
		totalPaidPrincipal := loanPrincipalAmountPaid + paymentPrincipalApplied
		if totalPaidPrincipal == loanPrincipalAmount && totalPaidFee == loanFeeAmount {
			isPaid = 1
		}

		// Update customer loan information
		customerLoan := models.CustomerLoan{
			ID:                 loanEntry.ID,
			TotalPaidPrincipal: null.FloatFrom(totalPaidPrincipal),
			TotalPaidFee:       null.FloatFrom(totalPaidFee),
			IsPaid:             null.IntFrom(int64(isPaid)),
			TotalPaidAmount:    null.FloatFrom(loanTotalPaidAmount + paymentPrincipalApplied + paymentFeeApplied),
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

		// Insert settlement information
		customerSettlement := models.CustomerSettlement{
			CustomerID:        null.IntFrom(int64(customer.ID)),
			CustomerLoanID:    null.IntFrom(int64(loanEntry.ID)),
			CustomerPaymentID: null.IntFrom(int64(customerPayment.ID)),
			SettlementAmount:  null.FloatFrom(paymentPrincipalApplied + paymentFeeApplied),
			PrincipalAmount:   null.FloatFrom(paymentPrincipalApplied),
			FeeAmount:         null.FloatFrom(paymentFeeApplied),
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
		"{amount}", strconv.FormatFloat(paymentAmount, 'f', -1, 64),
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
