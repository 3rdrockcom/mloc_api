package models

import null "gopkg.in/guregu/null.v3"

// Countries is an array of Country entries
type Countries []Country

// Country contains information about a country
type Country struct {
	ID               int         `db:"country_id" json:"country_id"` // CountryID should be match from database
	Name             string      `json:"name"`
	IsoCode2         string      `db:"iso_code_2" json:"iso_code_2"`
	IsoCode3         string      `db:"iso_code_3" json:"iso_code_3"`
	AddressFormat    string      `json:"address_format"`
	PostCodeRequired int         `json:"postcode_required"`
	Status           int         `json:"status"`
	MobilePrefix     null.String `json:"mobile_prefix"`
}

// TableName gets the name of the database table
func (c Country) TableName() string {
	return "tblCountry"
}
