package collection

import (
	"time"

	"github.com/epointpayment/mloc_api_go/app/models"

	"github.com/shopspring/decimal"
)

type Payload struct {
	Amount   decimal.Decimal
	Request  Request
	Response Response
}

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

type Push struct {
	Method          string
	Customer        models.CustomerInfo
	CustomerPayment models.CustomerPayment
	Transaction     PushTransaction
}

type PushTransaction struct {
	ID          string
	Description string
	Date        time.Time
	Amount      decimal.Decimal
}
