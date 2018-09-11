package models

import (
	null "gopkg.in/guregu/null.v3"
)

// TransactionsHistory is array of transactions
type TransactionsHistory []Transaction

// TableName gets the name of the database table
func (t TransactionsHistory) TableName() string {
	return "view_transaction_history"
}

// Transaction contains information about a loan transaction
type Transaction struct {
	CustomerID null.Int    `db:"fk_customer_id" json:"fk_customer_id"`
	Mode       null.String `db:"mode" json:"mode"`
	Amount     null.Float  `db:"amount" json:"amount"`
	Type       string      `db:"t_type" json:"t_type"`
	Date       null.Time   `db:"t_date" json:"t_date"`
}
