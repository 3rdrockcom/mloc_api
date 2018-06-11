package models

import (
	"gopkg.in/guregu/null.v3"
)

// Interest contains information about loan interest
type Interest struct {
	Active     null.String
	Percentage null.Float
	Fixed      null.Float
}

// TableName gets tblinterest from database
func (i Interest) TableName() string {
	return "tblInterest"
}
