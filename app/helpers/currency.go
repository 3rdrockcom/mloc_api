package helpers

import (
	"encoding/json"
	"errors"

	"github.com/labstack/gommon/log"
	"github.com/leekchan/accounting"
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

	// ErrCurrencyNotConfigured is an error given when the currency used does not have configuration entry
	ErrCurrencyNotConfigured = errors.New("Currency used is not configured")

	// ErrInvalidCurrencyAmount is an error given when the default currency is set to an invalid value
	ErrInvalidCurrencyAmount = errors.New("Invalid default currency amount used")
)

var currency = make(map[string]accounting.Accounting)

func init() {
	// Configure currency
	currency["USD"] = accounting.Accounting{Symbol: "$", Precision: 2, Thousand: ",", Decimal: "."}
	currency["MXN"] = accounting.Accounting{Symbol: "$", Precision: 2, Thousand: ",", Decimal: "."}

	// Set default currency information
	err := SetDefaultCurrency(DefaultCurrency)
	if err != nil {
		log.Fatal(err)
		return
	}
}

// SetDefaultCurrency sets the default currency used by the application
func SetDefaultCurrency(currencyCode string) (err error) {
	var currencyID, currencyPrecision int

	currencyID, currencyPrecision = iso4217.ByName(currencyCode)
	if currencyID == 0 {
		err = ErrInvalidCurrency
		return
	}

	if _, ok := currency[currencyCode]; !ok {
		err = ErrCurrencyNotConfigured
		return
	}

	DefaultCurrency = currencyCode
	DefaultCurrencyPrecision = int32(currencyPrecision)

	return
}

// FormatCurrency formats an amount using a specified currency format
func FormatCurrency(amount decimal.Decimal, currencyCode string) (formattedAmount string, err error) {
	c, ok := currency[currencyCode]
	if !ok {
		err = ErrCurrencyNotConfigured
		return
	}

	formattedAmount = c.FormatMoneyDecimal(amount)
	return
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
