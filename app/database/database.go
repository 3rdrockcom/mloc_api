package database

import (
	"log"

	dbx "github.com/go-ozzo/ozzo-dbx"
	_ "github.com/go-sql-driver/mysql" // Mysql driver for database library
)

// db is the database handler
var db *dbx.DB

// Database manages the database instances
type Database struct {
	driverName string
	dsn        string
	DB         *dbx.DB
}

// NewDatabase creates an instance of the service
func NewDatabase(driverName, dsn string) *Database {
	return &Database{
		driverName: driverName,
		dsn:        dsn,
	}
}

// Open creates a handler for the database
func (d Database) Open() error {
	var err error

	db, err = dbx.Open(d.driverName, d.dsn)
	if err != nil {
		return err
	}
	db.LogFunc = log.Printf

	return nil
}

// Close closes a handler for the database
func (d Database) Close() error {
	return d.DB.Close()
}

// GetInstance returns the database handler for this instance
func (d Database) GetInstance() *dbx.DB {
	return db
}

// Get returns the database handler
func Get() *dbx.DB {
	return db
}
