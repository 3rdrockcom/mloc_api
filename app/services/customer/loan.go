package customer

import (
	"database/sql"

	"github.com/epointpayment/mloc_api_go/app/models"
	dbx "github.com/go-ozzo/ozzo-dbx"
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
