package api

import (
	"database/sql"

	"github.com/epointpayment/mloc_api_go/app/models"

	dbx "github.com/go-ozzo/ozzo-dbx"
)

// DB is the database handler
var DB *dbx.DB

// APIService is a service that manages the API access
type APIService struct{}

// New creates an instance of the service
func New() *APIService {
	return &APIService{}
}

// DoAuth checks if a user's auth credentials are valid
func (as *APIService) DoAuth(username, password string) (isValid bool, err error) {
	authorizedUsers := make(map[string]string)

	authorizedUsers["EPOINT"] = "eyslTSh53q"

	if val, ok := authorizedUsers[username]; ok && val == password {
		isValid = true
	}

	return
}

// GetLoginKey gets the login API key
func (as *APIService) GetLoginKey() (entry *models.APIKey, err error) {
	entry = new(models.APIKey)

	err = DB.Select().
		Where(dbx.HashExp{"key": "LOGIN"}).
		One(entry)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrInvalidAPIKey
		}
		return nil, err
	}

	return
}

// GetRegistrationKey gets the registration API key
func (as *APIService) GetRegistrationKey() (entry *models.APIKey, err error) {
	entry = new(models.APIKey)

	err = DB.Select().
		Where(dbx.HashExp{"customer_id": 0}).
		AndWhere(dbx.NewExp("`key`!={:key}", dbx.Params{"key": "LOGIN"})).
		One(entry)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrInvalidAPIKey
		}
		return nil, err
	}

	return
}

// GetCustomerKey gets the customer API key
func (as *APIService) GetCustomerKey(key string) (entry *models.APIKey, err error) {
	entry, err = as.GetKey(key)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrInvalidAPIKey
		}
		return nil, err
	}

	return
}

// GetKey gets an API key
func (as *APIService) GetKey(key string) (entry *models.APIKey, err error) {
	entry = new(models.APIKey)

	err = DB.Select().
		Where(dbx.HashExp{"key": key}).
		One(entry)

	return
}

// GetKeyByCustomerID gets an API key by customer ID
func (as *APIService) GetKeyByCustomerID(customerID int) (entry *models.APIKey, err error) {
	entry = new(models.APIKey)

	err = DB.Select().
		Where(dbx.HashExp{"customer_id": customerID}).
		One(entry)

	return
}

// GetCustomerByCustomerUniqueID gets an API key by customer unique ID
func (as *APIService) GetCustomerByCustomerUniqueID(customerUniqueID string) (customer *models.Customer, err error) {
	customer = new(models.Customer)

	err = DB.Select().
		Where(dbx.HashExp{"cust_unique_id": customerUniqueID}).
		One(customer)
	if err != nil {
		return nil, err
	}

	return
}

// GetCustomerAccessKey gets a customer unique ID and associated customer API key
func (as *APIService) GetCustomerAccessKey(programID int, programCustomerID int, programCustomerMobile string) (k Key, err error) {
	// Initialize key service
	customerKey, err := NewKey(programID, programCustomerID, programCustomerMobile)
	if err != nil {
		return
	}

	// Get customer access key and customer unique ID
	k, err = customerKey.GetCustomerKey()
	if err != nil {
		return
	}

	return
}

// GenerateCustomerAccessKey generates a customer unique ID and associated customer API key
func (as *APIService) GenerateCustomerAccessKey(programID int, programCustomerID int, programCustomerMobile string) (k Key, err error) {
	// Initialize key service
	customerKey, err := NewKey(programID, programCustomerID, programCustomerMobile)
	if err != nil {
		return
	}

	// Get customer access key and customer unique ID
	k, err = customerKey.GenerateCustomerKey()
	if err != nil {
		return
	}

	return
}
