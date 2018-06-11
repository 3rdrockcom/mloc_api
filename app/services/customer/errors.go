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

	// ErrCustomerIncompleteInfo is given if customerinformation is not complete
	ErrCustomerIncompleteInfo = errors.New("Please provide complete customer information")

	// ErrProblemOccured is given if it can't get data from database or can't covert input to string
	ErrProblemOccured = errors.New("Some problems occurred, please try again")
)
