package epoint

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/url"
	"path"
	"strconv"

	"github.com/epointpayment/mloc_api_go/app/config"

	httpclient "github.com/ddliu/go-httpclient"
)

const (
	// Resources
	resourceMerchant = "merchant"

	// Endpoints
	endpointMerchantLogin           = "merchant_login"
	endpointMerchantFundTransfer    = "fund_transfer"
	endpointMerchantCustomerBalance = "get_customer_balance"
	endpointMerchantLogout          = "account_logout"
)

// cfg caches the config
var cfg config.Epoint

// EpointService is a service that manages the Epoint Payment API
type EpointService struct {
	sessionID string
}

// New creates an instance of the epoint service
func New() (e *EpointService, err error) {
	// Initialize config
	if cfg == (config.Epoint{}) {
		cfg = config.Get().Epoint
	}

	e = &EpointService{}
	return
}

// Response contains information given by the service
type Response struct {
	Code    string          `json:"ResponseCode"`
	Message json.RawMessage `json:"ResponseMessage"`
}

// SuccessResponseMessage contains success information
type SuccessResponseMessage struct {
	Message string `json:"message"`
}

// ErrorResponseMessage contains error information
type ErrorResponseMessage struct {
	Message string `json:"message"`
}

// LoginResponseMessage contains login information
type LoginResponseMessage struct {
	SessionID              string  `json:"session_id"`
	SessionValidityMinutes float64 `json:"session_validity_minutes"`
}

// GetLogin creates a new session
func (e *EpointService) GetLogin() (res LoginResponseMessage, err error) {
	// Generate a valid url
	u, err := e.generateURL(resourceMerchant, endpointMerchantLogin)
	if err != nil {
		return
	}

	// Make request
	resp, err := httpclient.Get(u.String(), map[string]string{
		"P01": strconv.Itoa(int(cfg.MTID)),
		"P02": cfg.Username,
		"P03": cfg.Password,
	})
	if err != nil {
		return
	}

	// Read response
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return
	}

	// Parse response
	r := new(Response)
	if err = json.Unmarshal(bodyBytes, r); err != nil {
		return
	}

	// Determine results
	switch r.Code {
	case "0000":
		if err := json.Unmarshal(r.Message, &res); err != nil {
			return res, err
		}

		e.sessionID = res.SessionID
		return res, err
	}

	// Handle error
	msg := new(ErrorResponseMessage)
	if err = json.Unmarshal(r.Message, msg); err != nil {
		return
	}

	// Send error
	err = errors.New(msg.Message)
	return
}

// FundTransferRequest contains request information for fund transfers
type FundTransferRequest struct {
	Amount          float64
	ClientReference string
	Source          string
	Destination     string
	Description     string
	MobileNumber    string
}

// FundTransferResponseMessage contains information about a fund transfer
type FundTransferResponseMessage struct {
	ClientReference string          `json:"client_reference_number"`
	TransactionID   json.RawMessage `json:"epoint_transaction_id"`
	Amount          json.Number     `json:"amount"`
}

// FundTransferResponse contains information about a fund transfer
type FundTransferResponse struct {
	ClientReference string
	TransactionID   string
	Amount          float64
}

// GetFundTransfer performs a fund transfer
func (e *EpointService) GetFundTransfer(req FundTransferRequest) (res FundTransferResponse, err error) {
	// Generate a valid url
	u, err := e.generateURL(resourceMerchant, endpointMerchantFundTransfer)
	if err != nil {
		return
	}

	bodyBytes := []byte(`{"ResponseCode":"0000","ResponseMessage":{"client_reference_number":"21","epoint_transaction_id":3580,"amount":"1"}}`)

	// bodyBytes := []byte(`{"ResponseCode":"1000","ResponseMessage":{"message":"Failed Validation.","validationMessages":{"v0":"The field UserDescription must be a string with a maximum length of 30."}}}`)

	if config.IsProd() {
		// Make request
		resp := new(httpclient.Response)
		resp, err = httpclient.Get(u.String(), map[string]string{
			"P01": e.sessionID,
			"P02": strconv.FormatFloat(req.Amount, 'f', -1, 64),
			"P03": req.ClientReference, // must unique
			"P04": req.Source,          // P,S,F, or customer ID
			"P05": req.Destination,     // P,S,F, or customer ID
			"P06": req.Description,     // 25 chars max
			"P07": req.MobileNumber,    // required, only if P04 is a customer_id
		})
		if err != nil {
			return
		}

		// Read response
		bodyBytes, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return
		}
	}

	// Parse response
	r := new(Response)
	if err = json.Unmarshal(bodyBytes, r); err != nil {
		return
	}

	// Determine results
	switch r.Code {
	case "0000":
		msg := new(FundTransferResponseMessage)
		if err = json.Unmarshal(r.Message, msg); err != nil {
			return
		}

		res.ClientReference = msg.ClientReference
		res.TransactionID = string(msg.TransactionID)
		res.Amount, err = msg.Amount.Float64()
		return
	}

	// Handle error
	msg := new(ErrorResponseMessage)
	if err = json.Unmarshal(r.Message, msg); err != nil {
		return
	}

	// Send error
	err = errors.New(msg.Message)
	return
}

// CustomerBalanceRequest contains request information about a customer's balance
type CustomerBalanceRequest struct {
	CustomerID   int
	MobileNumber string
}

// CustomerBalanceResponseMessage  contains information about a customer's balance
type CustomerBalanceResponseMessage struct {
	AvailableBalance float64 `json:"available_balance"`
}

// GetCustomerBalance obtains a customer's available balance
func (e *EpointService) GetCustomerBalance(req CustomerBalanceRequest) (res CustomerBalanceResponseMessage, err error) {
	// Generate a valid url
	u, err := e.generateURL(resourceMerchant, endpointMerchantCustomerBalance)
	if err != nil {
		return
	}

	bodyBytes := []byte(`{"ResponseCode":"0000","ResponseMessage":{"available_balance":95.25}}`)

	if config.IsProd() {
		// Make request
		resp := new(httpclient.Response)
		resp, err = httpclient.Get(u.String(), map[string]string{
			"P01": e.sessionID,
			"P02": strconv.FormatInt(int64(req.CustomerID), 10),
			"P03": req.MobileNumber,
		})
		if err != nil {
			return
		}

		// Read response
		bodyBytes, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return
		}
	}

	// Parse response
	r := new(Response)
	if err = json.Unmarshal(bodyBytes, r); err != nil {
		return
	}

	// Determine results
	switch r.Code {
	case "0000":
		if err = json.Unmarshal(r.Message, &res); err != nil {
			return
		}
		return
	}

	// Handle error
	msg := new(ErrorResponseMessage)
	if err = json.Unmarshal(r.Message, msg); err != nil {
		return
	}

	// Send error
	err = errors.New(msg.Message)
	return
}

// GetLogout logs out a session
func (e *EpointService) GetLogout() (res SuccessResponseMessage, err error) {
	// Generate a valid url
	u, err := e.generateURL(resourceMerchant, endpointMerchantLogout)
	if err != nil {
		return
	}

	// Make request
	resp, err := httpclient.Get(u.String(), map[string]string{
		"P01": e.sessionID,
	})
	if err != nil {
		return
	}

	// Read response
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return
	}

	// Parse response
	r := new(Response)
	if err = json.Unmarshal(bodyBytes, r); err != nil {
		return
	}

	// Determine results
	switch r.Code {
	case "0000":
		if err = json.Unmarshal(r.Message, &res); err != nil {
			return
		}
		return
	}

	// Handle error
	msg := new(ErrorResponseMessage)
	if err = json.Unmarshal(r.Message, msg); err != nil {
		return
	}

	// Send error
	err = errors.New(msg.Message)
	return
}

// generateURL generates an api endpoint url
func (e *EpointService) generateURL(resource string, endpoint string) (u *url.URL, err error) {
	// Parse api base url
	u, err = url.ParseRequestURI(cfg.BaseURL)
	if err != nil {
		return
	}

	// Merge url components
	u.Path = path.Join(u.Path, resource, endpoint)
	return
}
