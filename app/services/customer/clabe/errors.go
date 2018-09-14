package clabe

import "errors"

var (
	// ErrCLABEInvalidLength is an error given when the CLABE account number is not 18 digits in length
	ErrCLABEInvalidLength = errors.New("CLABE must be exactly 18 digits long")

	// ErrCLABENotNumeric is an error given when the CLABE account number contains non-numeric values
	ErrCLABENotNumeric = errors.New("CLABE must contain numeric digits only")

	// ErrCLABEChecksumNotValid is an error given when the checksum does not compute the correct value
	ErrCLABEChecksumNotValid = errors.New("CLABE checksum is not valid")
)
