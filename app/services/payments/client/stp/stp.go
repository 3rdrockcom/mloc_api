package stp

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"path"
	"strconv"

	httpclient "github.com/ddliu/go-httpclient"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

const (
	// Resources
	resourceInstitutions = "stp/institutions.php"
	resourceWebServices  = "stpgw/webservices.php"

	// Endpoints
	endpointGenerateCLABE = "get_clabe"
	endpointValidateCLABE = "validate_clabe"
	endpointSTPOut        = "stp_out"
)

// Client is a service that manages the STP Payment API
type Client struct {
	cfg Config
}

// Config contains information required for client
type Config struct {
	BaseURL   string
	ProgramID int64
	Username  string
	Password  string
}

// Validate checks if the configuration is invalid
func (m Config) Validate() (err error) {
	return validation.ValidateStruct(&m,
		validation.Field(&m.BaseURL, validation.Required, is.URL),
		validation.Field(&m.ProgramID, validation.Required, validation.Min(0)),
		validation.Field(&m.Username, validation.Required),
		validation.Field(&m.Password, validation.Required),
	)
}

// New creates an instance of the STP service
func New(cfg Config) (c *Client, err error) {
	// Validate config
	err = cfg.Validate()
	if err != nil {
		return
	}

	c = &Client{
		cfg: cfg,
	}
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

// GenerateCLABERequest contains request information for generating CLABE
type GenerateCLABERequest struct {
	ClientReference string
}

// GenerateCLABEResponseMessage contains response information for generating CLABE message
type GenerateCLABEResponseMessage struct {
	CLABE string `json:"ClabeNo"`
}

// GenerateCLABEResponse contains information about a CLABE
type GenerateCLABEResponse struct {
	CLABE string
}

// GenerateCLABE generates a CLABE number
func (c *Client) GenerateCLABE(req GenerateCLABERequest) (res GenerateCLABEResponse, err error) {
	// Generate a valid url
	u, err := c.generateURL(resourceWebServices)
	if err != nil {
		return
	}

	// Make request
	resp, err := httpclient.
		WithHeader("Authorization", generateBasicAuthHeader(c.cfg.Username, c.cfg.Password)).
		Post(u.String(), map[string]string{
			"method":  endpointGenerateCLABE,
			"prog_id": strconv.FormatInt(c.cfg.ProgramID, 10),
			"ref_no":  req.ClientReference,
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
		msg := new(GenerateCLABEResponseMessage)
		if err = json.Unmarshal(r.Message, msg); err != nil {
			return
		}

		res.CLABE = msg.CLABE
		return
	}

	// Handle error
	msg := new(ErrorResponseMessage)
	msg.Message = string(r.Message)

	// Send error
	err = errors.New(msg.Message)
	return
}

// ValidateCLABERequest contains request information for validating CLABE
type ValidateCLABERequest struct {
	CLABE string
}

// ValidateCLABEResponseMessage contains response information for validating CLABE message
type ValidateCLABEResponseMessage struct {
	CLABE           string      `json:"ClabeNo"`
	ProgramID       json.Number `json:"ProgID"`
	ClientReference string      `json:"RefNo"`
}

// ValidateCLABEResponse contains information about a CLABE
type ValidateCLABEResponse struct {
	CLABE           string
	ProgramID       int64
	ClientReference string
}

// ValidateCLABE validates a CLABE number
func (c *Client) ValidateCLABE(req ValidateCLABERequest) (res ValidateCLABEResponse, err error) {
	// Generate a valid url
	u, err := c.generateURL(resourceWebServices)
	if err != nil {
		return
	}

	// Make request
	resp, err := httpclient.
		WithHeader("Authorization", generateBasicAuthHeader(c.cfg.Username, c.cfg.Password)).
		Post(u.String(), map[string]string{
			"method":   endpointValidateCLABE,
			"clabe_no": req.CLABE,
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
		msg := new(ValidateCLABEResponseMessage)
		if err = json.Unmarshal(r.Message, msg); err != nil {
			return
		}

		res.CLABE = msg.CLABE
		res.ProgramID, _ = msg.ProgramID.Int64()
		res.ClientReference = msg.ClientReference
		return
	}

	// Handle error
	msg := new(ErrorResponseMessage)
	msg.Message = string(r.Message)

	// Send error
	err = errors.New(msg.Message)
	return
}

// GetInstitutionsResponseMessage contains response information about banking institutions
type GetInstitutionsResponseMessage []struct {
	ID   json.Number `json:"ID"`
	Name string      `json:"Institution"`
}

// Institution contains information about a banking institution
type Institution struct {
	ID   int64  `json:"ID"`
	Name string `json:"Institution"`
}

// GetInstitutionsResponse contains a list of banking institutions
type GetInstitutionsResponse struct {
	Institutions []Institution
}

// GetInstitutions gets a list of banking institutions for STP
func (c *Client) GetInstitutions() (res GetInstitutionsResponse, err error) {
	// Generate a valid url
	u, err := c.generateURL(resourceInstitutions)
	if err != nil {
		return
	}

	// Make request
	resp, err := httpclient.
		Get(u.String())
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
		msg := GetInstitutionsResponseMessage{}
		if err = json.Unmarshal(r.Message, &msg); err != nil {
			return
		}

		for _, institution := range msg {
			id, _ := institution.ID.Int64()
			res.Institutions = append(res.Institutions, Institution{
				ID:   id,
				Name: institution.Name,
			})
		}

		return
	}

	// Handle error
	msg := new(ErrorResponseMessage)
	msg.Message = string(r.Message)

	// Send error
	err = errors.New(msg.Message)
	return
}

// FundTransferOutboundRequest contains request information for STP disbursement
type FundTransferOutboundRequest struct {
	Amount      string
	Account     string
	Email       string
	Source      int64
	Destination int64
}

// FundTransferOutboundResponseMessage contains response information for STP disbursement
type FundTransferOutboundResponseMessage struct {
	TransactionID string `json:"TransId"`
}

// FundTransferOutboundResponse contains information about a STP disbursement
type FundTransferOutboundResponse struct {
	TransactionID string
}

// STPOut performs a STP disbursement
func (c *Client) STPOut(req FundTransferOutboundRequest) (res FundTransferOutboundResponse, err error) {
	// Generate a valid url
	u, err := c.generateURL(resourceWebServices)
	if err != nil {
		return
	}

	// Make request
	resp, err := httpclient.
		WithHeader("Authorization", generateBasicAuthHeader(c.cfg.Username, c.cfg.Password)).
		Post(u.String(), map[string]string{
			"method":      endpointSTPOut,
			"amount":      req.Amount,
			"bene_acct":   req.Account,
			"bene_email":  req.Email,
			"operante":    strconv.FormatInt(req.Source, 10),
			"contraparte": strconv.FormatInt(req.Destination, 10),
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
		msg := new(FundTransferOutboundResponseMessage)
		if err = json.Unmarshal(r.Message, msg); err != nil {
			return
		}

		res.TransactionID = msg.TransactionID
		return
	}

	// Handle error
	msg := new(ErrorResponseMessage)
	msg.Message = string(r.Message)

	// Send error
	err = errors.New(msg.Message)
	return
}

// generateURL generates an api resource url
func (c *Client) generateURL(resource string) (u *url.URL, err error) {
	// Parse api base url
	u, err = url.ParseRequestURI(c.cfg.BaseURL)
	if err != nil {
		return
	}

	// Merge url components
	u.Path = path.Join(u.Path, resource)
	return
}

// generateBasicAuthHeader encodes the username and password pair for basic auth
func generateBasicAuthHeader(username, password string) (auth string) {
	auth = base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	return fmt.Sprintf("Basic %s", auth)
}
