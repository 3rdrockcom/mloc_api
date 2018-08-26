package collection

import (
	"github.com/epointpayment/mloc_api_go/app/models"

	"github.com/shopspring/decimal"
)

type Request struct {
	Method          string
	Customer        models.CustomerInfo
	CustomerPayment models.CustomerPayment
	Description     string
}

type Response struct {
	ClientReference string
	TransactionID   string
	Amount          decimal.Decimal
}
