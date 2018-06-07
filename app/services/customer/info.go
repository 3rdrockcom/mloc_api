package customer

import (
	"database/sql"

	"github.com/epointpayment/mloc_api_go/app/models"

	dbx "github.com/go-ozzo/ozzo-dbx"
)

// Info manages customer information
type Info struct {
	cs *CustomerService
}

// Get gets basic customer information
func (i *Info) Get() (customer *models.Customer, err error) {
	customer = new(models.Customer)

	err = DB.Select().
		Where(dbx.HashExp{"id": i.cs.CustomerID}).
		One(customer)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrCustomerNotFound
		}
		return nil, err
	}

	return
}

// GetDetails gets detailed customer information
func (i *Info) GetDetails() (customerInfo *models.CustomerInfo, err error) {
	customerInfo = new(models.CustomerInfo)

	err = DB.Select().
		From("view_customer_info").
		Where(dbx.HashExp{"customer_id": i.cs.CustomerID}).
		One(customerInfo)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrCustomerNotFound
		}
		return nil, err
	}

	return
}

// Update updates customer information
func (i *Info) Update(customer *models.Customer) (err error) {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	customer.ID = i.cs.CustomerID
	err = tx.Model(customer).Update()
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
