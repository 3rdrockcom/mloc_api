package evault

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"

	httpclient "github.com/ddliu/go-httpclient"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

const (
	// Resources
	resourceDatastore = "datastore"

	// Endpoints
	endpointDatastoreStoreEntry = "store_entry"
	endpointDatastoreGetEntry   = "get_entry"
)

// Client is a service that manages the EVAULT API
type Client struct {
	cfg Config
}

// Config contains information required for client
type Config struct {
	BaseURL     string
	ProgramID   int64
	PartitionID int64
	Username    string
	Password    string
}

// Validate checks if the configuration is invalid
func (m Config) Validate() (err error) {
	return validation.ValidateStruct(&m,
		validation.Field(&m.BaseURL, validation.Required, is.URL),
		validation.Field(&m.ProgramID, validation.Required, validation.Min(0)),
		validation.Field(&m.PartitionID, validation.Required, validation.Min(0)),
		validation.Field(&m.Username, validation.Required),
		validation.Field(&m.Password, validation.Required),
	)
}

// New creates an instance of the evault service
func New(cfg Config) (e *Client, err error) {
	// Validate config
	err = cfg.Validate()
	if err != nil {
		return
	}

	e = &Client{
		cfg: cfg,
	}
	return
}

// Response contains information given by the service
type Response struct {
	Status       bool   `json:"status"`
	ResponseCode int    `json:"response_code"`
	Message      string `json:"message"`
}

// ErrorResponse contains error information given by the service
type ErrorResponse struct {
	Status       bool   `json:"status"`
	ResponseCode int    `json:"response_code"`
	Message      string `json:"error"`
}

// Response contains information given by the service
type EntryResponse struct {
	ID          int    `json:"id"`
	PartitionID int    `json:"partition_id"`
	Value       string `json:"value"`
}

func (c *Client) StoreEntry(data string) (res EntryResponse, err error) {
	// Generate a valid url
	u, err := c.generateURL(resourceDatastore, endpointDatastoreStoreEntry)
	if err != nil {
		return
	}

	// Set query string for url
	q := u.Query()
	q.Set("program_id", strconv.Itoa(int(c.cfg.ProgramID)))
	q.Set("partition_id", strconv.Itoa(int(c.cfg.PartitionID)))
	u.RawQuery = q.Encode()

	fmt.Println(u.String())

	// Make request
	resp, err := httpclient.
		WithHeader("Authorization", generateBasicAuthHeader(c.cfg.Username, c.cfg.Password)).
		Post(u.String(), map[string]string{
			"value": data,
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

	// Determine results
	switch resp.StatusCode {
	case http.StatusOK:
		if err = json.Unmarshal(bodyBytes, &res); err != nil {
			return res, err
		}

		return res, err
	}

	// Handle error
	r := new(ErrorResponse)
	if err = json.Unmarshal(bodyBytes, r); err != nil {
		return
	}

	// Send error
	err = errors.New(r.Message)
	return
}

func (c *Client) GetEntry(entryID int) (res EntryResponse, err error) {
	// Generate a valid url
	u, err := c.generateURL(resourceDatastore, endpointDatastoreGetEntry)
	if err != nil {
		return
	}

	// Make request
	resp, err := httpclient.
		WithHeader("Authorization", generateBasicAuthHeader(c.cfg.Username, c.cfg.Password)).
		Get(u.String(), map[string]string{
			"program_id":   strconv.Itoa(int(c.cfg.ProgramID)),
			"partition_id": strconv.Itoa(int(c.cfg.PartitionID)),
			"id":           strconv.Itoa(entryID),
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

	// Determine results
	switch resp.StatusCode {
	case http.StatusOK:
		if err = json.Unmarshal(bodyBytes, &res); err != nil {
			return res, err
		}

		return res, err
	}

	// Handle error
	r := new(ErrorResponse)
	if err = json.Unmarshal(bodyBytes, r); err != nil {
		return
	}

	// Send error
	err = errors.New(r.Message)
	return
}

// generateURL generates an api endpoint url
func (c *Client) generateURL(resource string, endpoint string) (u *url.URL, err error) {
	// Parse api base url
	u, err = url.ParseRequestURI(c.cfg.BaseURL)
	if err != nil {
		return
	}

	// Merge url components
	u.Path = path.Join(u.Path, resource, endpoint)
	return
}

// generateBasicAuthHeader encodes the username and password pair for basic auth
func generateBasicAuthHeader(username, password string) (auth string) {
	auth = base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	return fmt.Sprintf("Basic %s", auth)
}
