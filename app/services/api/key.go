package api

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/epointpayment/mloc_api_go/app/models"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/gommon/random"
	null "gopkg.in/guregu/null.v3"
	"gopkg.in/guregu/null.v3/zero"
)

// Key manages the customer API
type Key struct {
	programID             int
	programCustomerID     int
	programCustomerMobile string
	CustomerUniqueID      string   `json:"cust_unique_id"`
	ApiKey                string   `json:"api_key"`
	MLOCAccess            zero.Int `json:"mloc_access"`
	Registration          zero.Int `json:"registration"`
	TermsAndConditions    zero.Int `json:"term_and_condition"`
}

// Validate checks if the values in the struct are valid
func (k Key) Validate() error {
	return validation.ValidateStruct(&k,
		validation.Field(&k.programID, validation.Required),
		validation.Field(&k.programCustomerID, validation.Required),
		validation.Field(&k.programCustomerMobile, validation.Required),
	)
}

// NewKey creates an instance of the customer key service
func NewKey(programID int, programCustomerID int, programCustomerMobile string) (k *Key, err error) {
	k = &Key{
		programID:             programID,
		programCustomerID:     programCustomerID,
		programCustomerMobile: programCustomerMobile,
	}

	err = k.Validate()
	return
}

// GetCustomerKey gets a customer and associated customer API key
func (k *Key) GetCustomerKey() (customerKey Key, err error) {
	customerUniqueID := k.generateCustomerUniqueID()

	as := New()

	customer, err := as.GetCustomerInfoByCustomerUniqueID(customerUniqueID)
	if err == sql.ErrNoRows {
		entry, err := as.GetRegistrationKey()
		if err != nil {
			return customerKey, err
		}
		customerKey.ApiKey = entry.Key

		return customerKey, nil
	} else if err != nil {
		return
	}

	entry, err := as.GetKeyByCustomerID(customer.ID)
	if err != nil {
		return
	}

	customerKey.CustomerUniqueID = k.generateCustomerUniqueID()
	customerKey.ApiKey = entry.Key
	customerKey.MLOCAccess = customer.MLOCAccess
	customerKey.Registration = customer.Registration
	customerKey.TermsAndConditions = customer.TermsAndConditions

	return
}

// GenerateCustomerKey generates a customer and associated customer API key
func (k *Key) GenerateCustomerKey() (customerKey Key, err error) {
	customerUniqueID := k.generateCustomerUniqueID()

	as := New()

	_, err = as.GetCustomerInfoByCustomerUniqueID(customerUniqueID)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	if err == nil {
		err = ErrCustomerExists
		return
	}

	tx, err := DB.Begin()
	if err != nil {
		return
	}

	customer := &models.Customer{
		ProgramID:             null.IntFrom(int64(k.programID)),
		ProgramCustomerID:     null.IntFrom(int64(k.programCustomerID)),
		ProgramCustomerMobile: null.StringFrom(k.programCustomerMobile),
		CustomerUniqueID:      null.StringFrom(customerUniqueID),
	}
	err = tx.Model(customer).Insert()
	if err != nil {
		tx.Rollback()
		return
	}

	entry := &models.APIKey{
		CustomerID:  &customer.ID,
		Key:         k.generateAPIKey(),
		DateCreated: time.Now().UTC(),
	}
	err = tx.Model(entry).Insert()
	if err != nil {
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return
	}

	customerKey.ApiKey = entry.Key
	customerKey.CustomerUniqueID = customer.CustomerUniqueID.String

	return
}

// generateCustomerUniqueID generates an MD5 hash from customer program information
func (k *Key) generateCustomerUniqueID() string {
	str := strconv.Itoa(k.programID) + "_" + strconv.Itoa(k.programCustomerID) + "_" + k.programCustomerMobile

	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// generateAPIKey generates a random string of 32 characters
func (k *Key) generateAPIKey() string {
	return random.String(32)
}
