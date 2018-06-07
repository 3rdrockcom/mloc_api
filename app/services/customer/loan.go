package customer

import (
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

	// Clean up data
	for i := range transactionHistory {
		transactionHistory[i].Amount = transactionHistory[i].Amount.Truncate(2)
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

	// Clean up data
	for i := range transactionHistory {
		transactionHistory[i].Amount = transactionHistory[i].Amount.Truncate(2)
	}

	return
}
