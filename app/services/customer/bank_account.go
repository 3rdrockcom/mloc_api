package customer

import (
	"database/sql"
	"time"

	"github.com/epointpayment/mloc_api_go/app/config"
	"github.com/epointpayment/mloc_api_go/app/models"
	"github.com/epointpayment/mloc_api_go/app/services/customer/clabe"

	dbx "github.com/go-ozzo/ozzo-dbx"
)

const (
	// BankAccountTypeCLABE is a class name for CLABE type bank accounts
	BankAccountTypeCLABE = "CLABE"
)

// BankAccount manages customer bank accounts
type BankAccount struct {
	cs *CustomerService
}

// Get gets information about a bank account
func (a *BankAccount) Get(id int) (entry *models.CustomerBankAccount, err error) {
	entry = new(models.CustomerBankAccount)

	err = DB.Select().
		Where(dbx.HashExp{
			"id":          id,
			"customer_id": a.cs.CustomerID,
		}).
		One(entry)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrBankAccountNotFound
		}
		return nil, err
	}

	return
}

// GetUnmasked gets decrypted bank account information
func (a *BankAccount) GetUnmasked(id int) (entry *models.CustomerBankAccount, err error) {
	entry, err = a.Get(id)
	if err != nil {
		return
	}

	// Decrypt bank account number
	entry.AccountNumber, err = a.decrypt(entry.KmsID, entry.EvaultID)
	if err != nil {
		return
	}

	return
}

// GetAll gets information about all bank account
func (a *BankAccount) GetAll() (entries []models.CustomerBankAccount, err error) {
	err = DB.Select().
		From(models.CustomerBankAccount{}.TableName()).
		Where(dbx.HashExp{
			"customer_id": a.cs.CustomerID,
		}).
		All(&entries)
	if len(entries) == 0 {
		err = ErrBankAccountNotFound
	}
	if err != nil {
		return
	}

	return
}

// Create creates a bank account entry
func (a *BankAccount) Create(entry *models.CustomerBankAccount) (err error) {
	// Determine bank account type
	entry.AccountType, err = a.getType()
	if err != nil {
		return
	}

	// Check if bank account number is valid
	isValid, err := a.isValidAccountNumber(entry)
	if err != nil || isValid == false {
		err = ErrBankAccountNumberInvalid
		return
	}

	// Encrypt and store bank account number
	entry.KmsID, entry.EvaultID, err = a.encrypt(entry.AccountNumber)
	if err != nil {
		return
	}

	// Mask bank account number
	err = a.maskAccountNumber(entry)
	if err != nil {
		return
	}

	entry.CustomerID = a.cs.CustomerID
	entry.DateCreated = time.Now().UTC()

	tx, err := DB.Begin()
	if err != nil {
		return
	}

	err = tx.Model(entry).Insert()
	if err != nil {
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		err = tx.Rollback()
		return
	}

	return
}

// getType determines the bank account type (based on country of origin)
func (a *BankAccount) getType() (accountType string, err error) {
	// Get current country
	countryCode := config.Get().Country.Default

	// Determine bank account type based on country
	switch countryCode {
	case "MX":
		accountType = BankAccountTypeCLABE
	default:
		err = ErrBankAccountTypeUnknown
	}
	if err != nil {
		return
	}

	return
}

// isValidAccountNumber checks if the bank account number is valid
func (a *BankAccount) isValidAccountNumber(entry *models.CustomerBankAccount) (isValid bool, err error) {
	// Determine the bank account type and validate account number
	switch entry.AccountType {
	case BankAccountTypeCLABE:
		err = clabe.New(entry.AccountNumber).Validate()
	default:
		err = ErrBankAccountTypeUnknown
	}
	if err != nil {
		return false, err
	}

	isValid = true
	return
}

// maskAccountNumber will mask account number depending on the bank account type
func (a *BankAccount) maskAccountNumber(entry *models.CustomerBankAccount) (err error) {

	// Determine the bank account type and mask account number
	switch entry.AccountType {
	case BankAccountTypeCLABE:
		entry.AccountNumber = a.maskRange(entry.AccountNumber, 3, len(entry.AccountNumber)-4)
	default:
		err = ErrBankAccountTypeUnknown
	}
	if err != nil {
		return
	}

	return
}

// maskRange will obfuscate a range of characters within a string
func (a *BankAccount) maskRange(s string, start int, end int) string {
	rs := []rune(s)

	for i := 0; i < len(rs); i++ {
		if i >= start && i < end {
			rs[i] = '*'
		}
	}

	return string(rs)
}

// Update updates a bank account entry
func (a *BankAccount) Update(entry *models.CustomerBankAccount) (err error) {
	entry.CustomerID = a.cs.CustomerID
	entry.DateUpdated = time.Now().UTC()

	tx, err := DB.Begin()
	if err != nil {
		return
	}

	err = tx.Model(entry).Update(
		"Alias",
		"DateUpdated",
	)
	if err != nil {
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		err = tx.Rollback()
		return
	}

	return
}

// Delete deletes a bank account entry
func (a *BankAccount) Delete(entry *models.CustomerBankAccount) (err error) {
	entry.CustomerID = a.cs.CustomerID

	tx, err := DB.Begin()
	if err != nil {
		return
	}

	err = tx.Model(entry).Delete()
	if err != nil {
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		err = tx.Rollback()
		return
	}

	return
}
