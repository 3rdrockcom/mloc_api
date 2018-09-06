package disbursement

import (
	"github.com/epointpayment/mloc_api_go/app/models"

	"github.com/shopspring/decimal"
)

type Request struct {
	Method                  string
	Customer                models.CustomerInfo
	CustomerBankAccount     models.CustomerBankAccount
	CustomerLoanApplication models.CustomerLoanApplication
	Description             string
}

type Response struct {
	ClientReference string
	TransactionID   string
	Amount          decimal.Decimal
}
