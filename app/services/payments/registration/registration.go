package registration

import "github.com/epointpayment/mloc_api_go/app/models"

type Request struct {
	Method   string
	Customer models.CustomerInfo
}

type Response struct {
	Identifier string
}
