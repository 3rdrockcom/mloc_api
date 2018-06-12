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
	Code    string      `json:"ResponseCode"`
	Message interface{} `json:"ResponseMessage"`
}

// LoginResponse contains login information
type LoginResponse struct {
	SessionID              string
	SessionValidityMinutes float64
}

// GetLogin creates a new session
func (e *EpointService) GetLogin() (res LoginResponse, err error) {
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
	err = json.Unmarshal(bodyBytes, r)
	if err != nil {
		return
	}

	// Determine results
	if msg, ok := r.Message.(map[string]interface{}); ok {
		switch r.Code {
		case "0000":
			res.SessionID = msg["session_id"].(string)
			res.SessionValidityMinutes = msg["session_validity_minutes"].(float64)

			e.sessionID = res.SessionID
			return
		}
		err = errors.New(msg["message"].(string))
		return
	}

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

	// Make request
	resp, err := httpclient.Get(u.String(), map[string]string{
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
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return
	}

	// Parse response
	r := new(Response)
	err = json.Unmarshal(bodyBytes, r)
	if err != nil {
		return
	}

	// Determine results
	if msg, ok := r.Message.(map[string]interface{}); ok {
		switch r.Code {
		case "0000":
			res.ClientReference = msg["client_reference_number"].(string)
			res.TransactionID = msg["epoint_transaction_id"].(string)
			res.Amount = msg["amount"].(float64)
			return
		}
		err = errors.New(msg["message"].(string))
		return
	}

	return
}

// CustomerBalanceRequest contains request information about a customer's balance
type CustomerBalanceRequest struct {
	CustomerID   int
	MobileNumber string
}

// CustomerBalanceResponse contains information about a customer's balance
type CustomerBalanceResponse struct {
	AvailableBalance float64
}

// GetCustomerBalance obtains a customer's available balance
func (e *EpointService) GetCustomerBalance(req CustomerBalanceRequest) (res CustomerBalanceResponse, err error) {
	// Generate a valid url
	u, err := e.generateURL(resourceMerchant, endpointMerchantCustomerBalance)
	if err != nil {
		return
	}

	// Make request
	resp, err := httpclient.Get(u.String(), map[string]string{
		"P01": e.sessionID,
		"P02": strconv.FormatInt(int64(req.CustomerID), 10),
		"P03": req.MobileNumber,
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
	err = json.Unmarshal(bodyBytes, r)
	if err != nil {
		return
	}

	// Determine results
	if msg, ok := r.Message.(map[string]interface{}); ok {
		switch r.Code {
		case "0000":
			res.AvailableBalance = msg["available_balance"].(float64)
			return
		}
		err = errors.New(msg["message"].(string))
		return
	}

	return
}

// GetLogout logs out a session
func (e *EpointService) GetLogout() (err error) {
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
	err = json.Unmarshal(bodyBytes, r)
	if err != nil {
		return
	}

	// Determine results
	if msg, ok := r.Message.(map[string]interface{}); ok {
		switch r.Code {
		case "0000":
			return
		}
		err = errors.New(msg["message"].(string))
		return
	}

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
