package epoint

import (
	"strconv"

	"github.com/epointpayment/mloc_api_go/app/helpers"
	EPOINT "github.com/epointpayment/mloc_api_go/app/services/payments/client/epoint"
	"github.com/epointpayment/mloc_api_go/app/services/payments/collection"
	"github.com/epointpayment/mloc_api_go/app/services/payments/disbursement"
	"github.com/epointpayment/mloc_api_go/app/services/payments/driver"

	"github.com/shopspring/decimal"
)

type Driver struct{}

func New() (*Driver, error) {
	return &Driver{}, nil
}

func (d *Driver) Disbursement(req disbursement.Request) (res disbursement.Response, err error) {
	// Initialize epoint service
	es := new(EPOINT.EpointService)
	es, err = EPOINT.New()
	if err != nil {
		return
	}

	// Login to epoint service
	_, err = es.GetLogin()
	if err != nil {
		err = driver.ErrIssuerInvalidUserPassword
		return
	}

	// Transfer funds from prefund to customer wallet using epoint service
	fundTransferRequest := EPOINT.FundTransferRequest{
		Amount:          decimal.NewFromFloat(req.CustomerLoanApplication.LoanAmount.Float64).StringFixed(helpers.DefaultCurrencyPrecision),
		ClientReference: req.CustomerLoanApplication.ReferenceCode.String,
		Source:          "P",
		Destination:     strconv.FormatInt(req.Customer.ProgramCustomerID.Int64, 10),
		Description:     req.Description,
		MobileNumber:    req.Customer.MobileNumber.String,
	}
	ft := EPOINT.FundTransferResponse{}
	ft, err = es.GetFundTransfer(fundTransferRequest)
	if err != nil {
		err = driver.ErrIssuerFailedTransfer
		return
	}

	res = disbursement.Response{
		ClientReference: ft.ClientReference,
		TransactionID:   ft.TransactionID,
		Amount:          decimal.NewFromFloat(ft.Amount),
	}

	return
}

func (d *Driver) Collection(req collection.Request) (res collection.Response, err error) {
	// Initialize epoint service
	es := new(EPOINT.EpointService)
	es, err = EPOINT.New()
	if err != nil {
		return
	}

	// Login to epoint service
	_, err = es.GetLogin()
	if err != nil {
		err = driver.ErrIssuerInvalidUserPassword
		return
	}

	// Get customer balance information
	customerBalanceRequest := EPOINT.CustomerBalanceRequest{
		CustomerID:   int(req.Customer.ProgramCustomerID.Int64),
		MobileNumber: req.Customer.MobileNumber.String,
	}
	cb, err := es.GetCustomerBalance(customerBalanceRequest)
	if err != nil {
		err = driver.ErrIssuerUnableToAccessBalance
		return
	}

	// Check if there is enough funds available in wallet for payment
	if decimal.NewFromFloat(req.CustomerPayment.PaymentAmount.Float64).GreaterThan(cb.AvailableBalance) {
		err = driver.ErrIssuerInsufficientFunds
		return
	}

	// Transfer funds from customer wallet to settlement using epoint service
	fundTransferRequest := EPOINT.FundTransferRequest{
		Amount:          decimal.NewFromFloat(req.CustomerPayment.PaymentAmount.Float64).StringFixed(helpers.DefaultCurrencyPrecision),
		ClientReference: req.CustomerPayment.ReferenceCode.String,
		Source:          strconv.FormatInt(req.Customer.ProgramCustomerID.Int64, 10),
		Destination:     "S",
		Description:     req.Description,
		MobileNumber:    req.Customer.MobileNumber.String,
	}
	ft := EPOINT.FundTransferResponse{}
	ft, err = es.GetFundTransfer(fundTransferRequest)
	if err != nil {
		err = driver.ErrIssuerFailedTransfer
		return
	}

	res = collection.Response{
		ClientReference: ft.ClientReference,
		TransactionID:   ft.TransactionID,
		Amount:          decimal.NewFromFloat(ft.Amount),
	}

	return
}
