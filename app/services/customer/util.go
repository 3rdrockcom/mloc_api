package customer

import (
	"math"

	"github.com/labstack/gommon/random"
)

// generateRandomKey generates a random key
func generateRandomKey(length uint8) string {
	if length == 0 {
		length = 20
	}

	randomKey := random.New().String(length, random.Uppercase+random.Numeric)
	return randomKey
}

// numberFormat formats a number to a certain amount of digits
func numberFormat(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(int(num*output+math.Copysign(0.5, num*output))) / output
}
