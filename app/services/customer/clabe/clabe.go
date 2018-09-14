package clabe

import "strconv"

type CLABE struct {
	AccountNumber string
}

func New(accountNumber string) *CLABE {
	return &CLABE{AccountNumber: accountNumber}
}

func (c *CLABE) Validate() (err error) {
	clabe := c.AccountNumber

	// Check if CLABE is 18 digits long
	if len(clabe) != 18 {
		err = ErrCLABEInvalidLength
		return
	}

	// Check if CLABE is numeric
	_, err = strconv.Atoi(clabe)
	if err != nil {
		err = ErrCLABENotNumeric
		return
	}

	// Parse CLABE account number into components
	bankCode := clabe[0:3]
	branchOfficeCode := clabe[3:6]
	accountNumber := clabe[6:17]
	controlDigit := clabe[17:]

	digits := make([]int, 18)
	controlDigitIndex := len(digits) - 1

	// Calculate checksum
	weights := []int{3, 7, 1}
	for i, r := range bankCode + branchOfficeCode + accountNumber {
		digit, _ := strconv.Atoi(string(r))

		// Multiply each digit by a specified weight and take modulus 10
		digits[i] = (digit * weights[i%3]) % 10

		// Sum all of the calculated products
		digits[controlDigitIndex] += digits[i]
	}
	// Sum all of the calculated products, and take modulus 10 again
	digits[controlDigitIndex] = 10 - digits[controlDigitIndex]%10

	// Check if calculated control digit matches submitted control digit
	if strconv.Itoa(digits[controlDigitIndex]) != controlDigit {
		err = ErrCLABEChecksumNotValid
		return
	}

	return
}
