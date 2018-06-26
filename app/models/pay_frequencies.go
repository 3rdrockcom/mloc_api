package models

import null "gopkg.in/guregu/null.v3"

// PayFrequencies is an array of pay frequency entries
type PayFrequencies []PayFrequency

// PayFrequency contains information about a pay frequency entry
type PayFrequency struct {
	ID          int         `db:"id" json:"id"`
	Description null.String `db:"description" json:"description"`
}

// TableName gets the name of the database table
func (c PayFrequency) TableName() string {
	return "tblPayFrequency"
}
