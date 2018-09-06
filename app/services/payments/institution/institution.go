package institution

import (
	"github.com/epointpayment/mloc_api_go/app/models"
)

type Request struct {
	Method   string
	Customer models.CustomerInfo
}

type Response struct {
	Institutions []Institution `json:"institutions"`
}

type Institution struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
