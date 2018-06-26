package services

import (
	"github.com/epointpayment/mloc_api_go/app/database"
	"github.com/epointpayment/mloc_api_go/app/services/api"
	"github.com/epointpayment/mloc_api_go/app/services/customer"
	"github.com/epointpayment/mloc_api_go/app/services/lookup"

	dbx "github.com/go-ozzo/ozzo-dbx"
)

// db is the database handler
var db *dbx.DB

// Services boots application-specific services
type Services struct{}

// New starts the service setup process
func New(DB *database.Database) error {
	db = DB.GetInstance()

	// Attach the database handler to service
	api.DB = db
	customer.DB = db
	lookup.DB = db

	return nil
}
