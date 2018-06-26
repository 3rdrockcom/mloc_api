package helpers

import (
	"encoding/json"
	"errors"

	"github.com/labstack/gommon/log"
	"github.com/rmg/iso4217"
	"github.com/shopspring/decimal"
	null "gopkg.in/guregu/null.v3"
)

var (
	// DefaultCurrency is the alphabetic code used to represent the default currency
	DefaultCurrency = "USD"

	// DefaultCurrencyPrecision is the number of digits to keep past the decimal point
	DefaultCurrencyPrecision int32

	// ErrInvalidCurrency is an error given when the default currency is set to an invalid value
	ErrInvalidCurrency = errors.New("Invalid default currency used")

	// ErrInvalidCurrencyAmount is an error given when the default currency is set to an invalid value
	ErrInvalidCurrencyAmount = errors.New("Invalid default currency amount used")
)

func init() {
	var currencyCode, currencyPrecision int

	currencyCode, currencyPrecision = iso4217.ByName(DefaultCurrency)
	if currencyCode == 0 {
		log.Fatal(ErrInvalidCurrency)
	}

	DefaultCurrencyPrecision = int32(currencyPrecision)
}

// ValidateCurrencyAmount is used by the validator to check if value is valid
func ValidateCurrencyAmount(value interface{}) error {
	var s string

	switch value.(type) {
	case null.String:
		s = value.(null.String).String
	case json.Number:
		s = value.(json.Number).String()
	}

	dec, err := decimal.NewFromString(s)
	if err != nil {
		return err
	}

	if dec.Truncate(DefaultCurrencyPrecision).Equal(dec) {
		if dec.GreaterThan(decimal.Zero) {
			return nil
		}
	}

	return ErrInvalidCurrencyAmount
}
