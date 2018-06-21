package customer

import (
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
