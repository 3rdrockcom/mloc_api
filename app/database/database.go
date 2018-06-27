package database

import (
	"fmt"
	"time"

	"github.com/epointpayment/mloc_api_go/app/log"

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

	// setup logging
	db.PerfFunc = logFunc

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

// logFunc is used to log the SQL execution time.
func logFunc(ns int64, sql string, execute bool) {
	log.GetInstance().WithFields(map[string]interface{}{
		"elasped":       ns,
		"elasped_human": fmt.Sprintf("%fms", float64(ns)/1000000.0),
		"sql":           sql,
		"execute":       execute,
	}).Debug("SQL Statement")
}
