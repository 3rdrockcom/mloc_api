package customer

import (
	"database/sql"
	"time"

	"github.com/epointpayment/mloc_api_go/app/models"

	dbx "github.com/go-ozzo/ozzo-dbx"
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
func (a *BankAccount) Create(entry *models.CustomerBankAccount) ( /*entry *models.CustomerBankAccount,*/ err error) {
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
