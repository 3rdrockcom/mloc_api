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

	// ErrCustomerIncompleteInfo is given if customerinformation is not complete
	ErrCustomerIncompleteInfo = errors.New("Please provide complete customer information")

	// ErrProblemOccured is given if it can't get data from database or can't covert input to string
	ErrProblemOccured = errors.New("Some problems occurred, please try again")
)
