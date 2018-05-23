package models

import "time"

// APIKeys is an array of APIKey entries
type APIKeys []APIKey

// APIKey contains information about a key
type APIKey struct {
	ID          int
	CustomerID  int
	Key         string
	DateCreated time.Time
}

// TableName gets the name of the database table
func (a APIKey) TableName() string {
	return "api_keys"
}
