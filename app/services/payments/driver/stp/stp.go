package stp

import (
	"strconv"

	"github.com/epointpayment/mloc_api_go/app/config"
	"github.com/epointpayment/mloc_api_go/app/helpers"
	STP "github.com/epointpayment/mloc_api_go/app/services/payments/client/stp"
	"github.com/epointpayment/mloc_api_go/app/services/payments/collection"
	"github.com/epointpayment/mloc_api_go/app/services/payments/disbursement"
	"github.com/epointpayment/mloc_api_go/app/services/payments/driver"
	"github.com/epointpayment/mloc_api_go/app/services/payments/registration"

	"github.com/shopspring/decimal"
)

// Driver manages the adapter for client
type Driver struct{}

// New creates a new instance of the client adapter
func New() (*Driver, error) {
	return &Driver{}, nil
}

func (d *Driver) Register(req registration.Request) (res registration.Response, err error) {
	// Get config for stp service
	cfg, err := d.getConfig()
	if err != nil {
		err = driver.ErrIssuerInvalidConfig
		return
	}
	cfg.ProgramID = req.Customer.ProgramCustomerID.Int64

	// Initialize stp service
	client := new(STP.Client)
	client, err = STP.New(cfg)
	if err != nil {
		return
	}

	// Generate CLABE
	reg, err := client.GenerateCLABE(STP.GenerateCLABERequest{
		ClientReference: req.Customer.CustUniqueID.String,
	})
	if err != nil {
		return
	}

	res = registration.Response{
		Identifier: reg.CLABE,
	}

	return
}

// Disbursement is pay out of funds (loan proceeds) to the borrower
func (d *Driver) Disbursement(req disbursement.Request) (res disbursement.Response, err error) {
	// Get config for stp service
	cfg, err := d.getConfig()
	if err != nil {
		err = driver.ErrIssuerInvalidConfig
		return
	}
	cfg.ProgramID = req.Customer.ProgramCustomerID.Int64

	// Initialize stp service
	client := new(STP.Client)
	client, err = STP.New(cfg)
	if err != nil {
		return
	}

	bankCode, err := strconv.Atoi(req.CustomerBankAccount.BankCode)
	if err != nil {
		return
	}

	// Transfer funds from STP to bank account
	fundTransferRequest := STP.FundTransferOutboundRequest{
		Amount: decimal.NewFromFloat(req.CustomerLoanApplication.LoanAmount.Float64).StringFixed(helpers.DefaultCurrencyPrecision),
		// Account:     "846180000400000001",
		Account: req.CustomerBankAccount.AccountNumber,
		Email:   req.Customer.Email.String,
		Source:  90646,
		// Destination: 846,
		Destination: int64(bankCode),
	}
	ft, err := client.STPOut(fundTransferRequest)
	if err != nil {
		err = driver.ErrIssuerFailedTransfer
		return
	}

	res = disbursement.Response{
		TransactionID: ft.TransactionID,
		Amount:        decimal.NewFromFloat(req.CustomerLoanApplication.LoanAmount.Float64),
	}

	return
}

// Collection is pay out of funds (loan proceeds) to the lender
func (d *Driver) Collection(req collection.Request) (res collection.Response, err error) {
	return
}

// getConfig gets the client configuration
func (d *Driver) getConfig() (c STP.Config, err error) {
	// Get config for epoint service
	cfg := config.Get().STP

	c = STP.Config{
		BaseURL: cfg.BaseURL,
		// ProgramID int64
		Username: cfg.Username,
		Password: cfg.Password,
	}

	return
}
