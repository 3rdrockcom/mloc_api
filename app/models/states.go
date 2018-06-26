package models

import null "gopkg.in/guregu/null.v3"

// States is an array of State entries
type States []State

// State contains information about a state
type State struct {
	ID        int      `db:"state_id" json:"state_id"`
	Name      string   `db:"state" json:"state"`
	Code      string   `db:"state_code" json:"state_code"`
	CountryID null.Int `db:"country_id" json:"country_id"`
}

// TableName gets the name of the database table
func (c State) TableName() string {
	return "tblState"
}
