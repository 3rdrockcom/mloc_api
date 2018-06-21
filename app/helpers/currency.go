package helpers

import (
	"errors"

	"github.com/labstack/gommon/log"
	"github.com/rmg/iso4217"
)

var (
	// DefaultCurrency is the alphabetic code used to represent the default currency
	DefaultCurrency = "USD"

	// DefaultCurrencyPrecision is the number of digits to keep past the decimal point
	DefaultCurrencyPrecision int32

	// ErrInvalidCurrency is an error given when the default currency is set to an invalid value
	ErrInvalidCurrency = errors.New("Invalid default currency used")
)

func init() {
	var currencyCode, currencyPrecision int

	currencyCode, currencyPrecision = iso4217.ByName(DefaultCurrency)
	if currencyCode == 0 {
		log.Fatal(ErrInvalidCurrency)
	}

	DefaultCurrencyPrecision = int32(currencyPrecision)
}
