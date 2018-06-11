package models

import null "gopkg.in/guregu/null.v3"

// SystemSettings contains an array of system settings
type SystemSettings []SystemSetting

// SystemSetting contains information about a global setting
type SystemSetting struct {
	ID           int
	Name         null.String
	Code         null.String
	Description  null.String
	Value        null.String
	SettingType  null.String
	IsActive     null.String
	SMSMessage   null.String `db:"sms_message"`
	EmailMessage null.String
	Subject      null.String
	From         null.String
	To           null.String
	Cc           null.String
	Bcc          null.String
	UpdatedBy    null.Int
}

// TableName gets the name of the database table
func (ss SystemSetting) TableName() string {
	return "tblSystemSettings"
}
