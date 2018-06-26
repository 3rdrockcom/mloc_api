package database

import (
	"log"
	"time"

	dbx "github.com/go-ozzo/ozzo-dbx"
	_ "github.com/go-sql-driver/mysql" // Mysql driver for database library
)

// db is the database handler
var db *dbx.DB

var (
	// maxOpenConns sets the maximum number of open connections to the database.
	maxOpenConns = 50

	// maxIdleConns sets the maximum number of connections in the idle connection pool.
	maxIdleConns = 0

	// connMaxLifetime sets the maximum amount of time a connection may be reused.
	connMaxLifetime = 5 * time.Second
)

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

	// Set database connection settings
	db.DB().SetMaxOpenConns(maxOpenConns)
	db.DB().SetMaxIdleConns(maxIdleConns)
	db.DB().SetConnMaxLifetime(connMaxLifetime)

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
