package models

// Cities is an array of city entries
type Cities []City

// City contains information about a CityID
type City struct {
	ID        int    `db:"city_id" json:"city_id"`
	Name      string `db:"city" json:"city"`
	StateCode string `db:"state_code" json:"state_code"`
}

// TableName gets the name of the database table
func (c City) TableName() string {
	return "tblCity"
}
