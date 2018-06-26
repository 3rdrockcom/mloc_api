package models

import null "gopkg.in/guregu/null.v3"

// IncomeSources is an array of IncomeSource entries
type IncomeSources []IncomeSource

// IncomeSource contains information about an income source
type IncomeSource struct {
	ID          int         `db:"id" json:"id"`
	Description null.String `db:"description" json:"description"`
}

// TableName gets the name of the database table
func (c IncomeSource) TableName() string {
	return "tblIncomeSource"
}
