package driver

import (
	"errors"
)

var (
	// ErrIssuerInvalidUserPassword is given if it username or password is wrong
	ErrIssuerInvalidUserPassword = errors.New("Invalid issuer username or password")

	// ErrIssuerUnableToAccessBalance is given when the system is unable to obtain the customer's balance
	ErrIssuerUnableToAccessBalance = errors.New("Cannot retrieve customer balance from issuer")

	// ErrIssuerInsufficientFunds is given when the customer does not have enought funds in their wallet
	ErrIssuerInsufficientFunds = errors.New("You dont have enough available funds in your wallet")

	// ErrIssuerFailedTransfer is given if an amount cannot transfer to/from wallet
	ErrIssuerFailedTransfer = errors.New("Cannot transfer amount with issuer: Please contact administrator")
)
