package router

import (
	"net/http"

	"github.com/epointpayment/mloc_api_go/app/controllers"
	API "github.com/epointpayment/mloc_api_go/app/services/api"
	Customer "github.com/epointpayment/mloc_api_go/app/services/customer"
	Lookup "github.com/epointpayment/mloc_api_go/app/services/lookup"

	"github.com/labstack/echo"
)

// appendErrorHandler handles errors for the router
func (r *Router) appendErrorHandler() {
	r.e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		message := err.Error()
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			message = he.Message.(string)
		}

		// Override status code based on error responses
		switch message {
		// API Service
		case API.ErrInvalidAPIKey.Error():
			code = http.StatusForbidden
		case API.ErrInvalidProgramCustomerID.Error():
			code = http.StatusBadRequest

		// Customer Service
		case Customer.ErrInvalidUniqueCustomerID.Error():
			code = http.StatusForbidden
		case Customer.ErrCustomerNotFound.Error():
			code = http.StatusNotFound
		case Customer.ErrTransactionNotFound.Error():
			code = http.StatusNotFound
		case Customer.ErrInvalidLoanAmount.Error():
			code = http.StatusBadRequest
		case Customer.ErrLoanNotFound.Error():
			code = http.StatusNotFound
		case Customer.ErrNotEnoughAvailableCredit.Error():
			code = http.StatusBadRequest
		case Customer.ErrCustomerIncompleteInfo.Error():
			code = http.StatusBadRequest
		case Customer.ErrProblemOccured.Error():
			code = http.StatusBadRequest
		case Customer.ErrLoanCreditLimitNotFound.Error():
			code = http.StatusBadRequest
		case Customer.ErrLoanInterestNotFound.Error():
			code = http.StatusBadRequest
		case Customer.ErrLoanFeeNotFound.Error():
			code = http.StatusBadRequest
		case Customer.ErrProcessLoanApplication.Error():
			code = http.StatusBadRequest
		case Customer.ErrProcessLoanPayment.Error():
			code = http.StatusBadRequest
		case
			Customer.ErrInvalidNextPayDate.Error(),
			Customer.ErrInvalidFollowingPayDate.Error():
			code = http.StatusBadRequest

		// Issuer / Epoint Service
		case Customer.ErrIssuerInvalidUserPassword.Error():
			code = http.StatusBadRequest
		case Customer.ErrIssuerUnableToAccessBalance.Error():
			code = http.StatusBadRequest
		case Customer.ErrIssuerInsufficientFunds.Error():
			code = http.StatusBadRequest
		case Customer.ErrIssuerFailedTransfer.Error():
			code = http.StatusBadRequest

		// Lookup Service
		case Lookup.ErrInvalidCountryID.Error():
			code = http.StatusBadRequest
		case Lookup.ErrCountryNotFound.Error():
			code = http.StatusNotFound
		case Lookup.ErrInvalidStateID.Error():
			code = http.StatusBadRequest
		case Lookup.ErrStateNotFound.Error():
			code = http.StatusNotFound
		case Lookup.ErrInvalidStateCode.Error():
			code = http.StatusBadRequest
		case Lookup.ErrInvalidCityID.Error():
			code = http.StatusBadRequest
		case Lookup.ErrCityNotFound.Error():
			code = http.StatusNotFound
		case Lookup.ErrInvalidIncomeSourceID.Error():
			code = http.StatusBadRequest
		case Lookup.ErrIncomeSourceNotFound.Error():
			code = http.StatusNotFound
		case Lookup.ErrInvalidIPayFrequencyID.Error():
			code = http.StatusBadRequest
		case Lookup.ErrPayFrequencyNotFound.Error():
			code = http.StatusNotFound

		// Unknown error
		default:
			if _, ok := err.(*echo.HTTPError); !ok {
				message = "Internal Error"
			}
		}

		// Log errors
		c.Logger().Error(err)

		// Send response
		if !c.Response().Committed {
			if c.Request().Method == echo.HEAD { // Issue #608
				err = c.NoContent(code)
			} else {
				// Send error in a specific format
				err = controllers.SendErrorResponse(c, code, message)
			}
			if err != nil {
				r.e.Logger.Error(err)
			}
		}
	}
}
