package auth

import (
	API "github.com/epointpayment/mloc_api_go/app/services/api"

	"github.com/labstack/echo"
)

// BasicValidator is a validator used for basic auth middleware
func BasicValidator(username, password string, roles []string, c echo.Context) (isValid bool, err error) {
	// Initialize API service
	sa := API.New()

	// Check is user is authorized
	for _, role := range roles {
		switch role {
		case "system":
			isValid, err = sa.DoSystemAuth(username, password)
		case "default":
			isValid, err = sa.DoAuth(username, password)
		}
		if err != nil || isValid == true {
			return
		}

	}

	return
}
