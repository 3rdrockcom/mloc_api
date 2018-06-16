package models

import (
	"gopkg.in/guregu/null.v3"
)

// Fee contains information about a loan fee
type Fee struct {
	Active     null.String
	Percentage null.Float
	Fixed      null.Float
}

// TableName gets the name of the database table
func (f Fee) TableName() string {
	return "tblFee"
}
