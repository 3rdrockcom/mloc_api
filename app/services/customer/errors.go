package customer

import "errors"

var (
	// ErrInvalidUniqueCustomerID is an error shown when customer unique ID is not a valid
	ErrInvalidUniqueCustomerID = errors.New("Invalid Unique Customer ID")

	// ErrCustomerNotFound is an error for a non-existent customer
	ErrCustomerNotFound = errors.New("Customer not found")

	// ErrTransactionNotFound is an error given when no transactions were found
	ErrTransactionNotFound = errors.New("No transaction(s) were found")

	// ErrLoanNotFound is an error given when a customer loan is not found
	ErrLoanNotFound = errors.New("No customer loan was found")

	// ErrLoanCreditLimitNotFound is an error given when a loan credit limit entry is not found
	ErrLoanCreditLimitNotFound = errors.New("No loan credit limit entry was found")

	// ErrLoanInterestNotFound is an error given when a loan interest entry is not found
	ErrLoanInterestNotFound = errors.New("No loan interest entry was found")

	// ErrLoanFeeNotFound is an error given when a loan fee entry is not found
	ErrLoanFeeNotFound = errors.New("No loan fee entry was found")

	// ErrInvalidLoanAmount is an error given when a loan amount is not valid
	ErrInvalidLoanAmount = errors.New("Invalid loan amount")

	// ErrNotEnoughAvailableCredit is given when loan amount is larger than the available credit
	ErrNotEnoughAvailableCredit = errors.New("You dont have enough available credit")

	// ErrProcessLoanApplication is given when a loan application fails to process
	ErrProcessLoanApplication = errors.New("Error while processing loan application")

	// ErrProcessLoanPayment is given when a loan payment fails to process
	ErrProcessLoanPayment = errors.New("Error while processing loan payment")

	// ErrIssuerInvalidUserPassword is given if it username or password is wrong
	ErrIssuerInvalidUserPassword = errors.New("Invalid issuer username or password")

	// ErrIssuerUnableToAccessBalance is given when the system is unable to obtain the customer's balance
	ErrIssuerUnableToAccessBalance = errors.New("Cannot retrieve customer balance from issuer")

	// ErrIssuerInsufficientFunds is given when the customer does not have enought funds in their wallet
	ErrIssuerInsufficientFunds = errors.New("You dont have enough available funds in your wallet")

	// ErrIssuerFailedTransfer is given if an amount cannot transfer to/from wallet
	ErrIssuerFailedTransfer = errors.New("Cannot transfer amount with issuer: Please contact administrator")

	// ErrCustomerIncompleteInfo is given if customerinformation is not complete
	ErrCustomerIncompleteInfo = errors.New("Please provide complete customer information")

	// ErrProblemOccured is given if it can't get data from database or can't covert input to string
	ErrProblemOccured = errors.New("Some problems occurred, please try again")
)
