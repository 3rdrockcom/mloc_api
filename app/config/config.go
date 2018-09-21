package config

import (
	"fmt"
	"strconv"

	"github.com/gobuffalo/envy"
)

// Configuration contains the application configuration
type Configuration struct {
	Environment string
	Application Application
	Country     Country
	Currency    Currency
	Server      Server
	DB          Database
	Mail        Mail
	SMS         SMS
	Epoint      Epoint
	STP         STP
	KMS         KMS
	Evault      Evault
}

// Application contains application information
type Application struct {
	Name    string
	Build   string
	Version string
}

// Country contains country information
type Country struct {
	Default string
}

// Currency contains currency information
type Currency struct {
	Default string
}

// Server contains server information
type Server struct {
	Host string
	Port int64
}

// Database contains database information
type Database struct {
	Driver   string
	Host     string
	Port     int64
	Database string
	Username string
	Password string
	Flags    string
	DSN      string
}

// Mail contains mail information
type Mail struct {
	Driver      string
	Host        string
	Port        int64
	Encryption  string
	Username    string
	Password    string
	FromName    string
	FromAddress string
	ToAddress   string
}

// SMS contains sms information
type SMS struct {
	Driver   string
	Username string
	Password string
	FromName string
	ToMobile string
}

// Epoint contains epoint merchant api information
type Epoint struct {
	BaseURL  string
	MTID     int64
	Username string
	Password string
}

// STP contains stp api information
type STP struct {
	BaseURL  string
	Username string
	Password string
}

// KMS contains kms api information
type KMS struct {
	BaseURL   string
	ProgramID int64
	Username  string
	Password  string
}

// Evault contains evault api information
type Evault struct {
	BaseURL     string
	ProgramID   int64
	PartitionID int64
	Username    string
	Password    string
}

// cfg contains the processed configuration values
var cfg Configuration

// New processes the configuration values
func New() (Configuration, error) {
	// Base
	cfg.Environment = envy.Get("ENVIRONMENT", "development")

	// Application
	cfg.Application.Name = envy.Get("NAME", "app")
	cfg.Application.Build = Build
	cfg.Application.Version = Version

	// Country
	cfg.Country.Default = envy.Get("COUNTRY_DEFAULT", "US")

	// Currency
	cfg.Currency.Default = envy.Get("CURRENCY_DEFAULT", "USD")

	// Server
	cfg.Server.Host = envy.Get("HOST", "localhost")
	cfg.Server.Port, _ = strconv.ParseInt(envy.Get("PORT", "3000"), 10, 64)

	// Database
	cfg.DB.Driver, _ = envy.MustGet("DB_CONNECTION")
	cfg.DB.Host = envy.Get("DB_HOST", "localhost")
	cfg.DB.Port, _ = strconv.ParseInt(envy.Get("DB_PORT", "3306"), 10, 64)
	cfg.DB.Database, _ = envy.MustGet("DB_DATABASE")
	cfg.DB.Username, _ = envy.MustGet("DB_USERNAME")
	cfg.DB.Password, _ = envy.MustGet("DB_PASSWORD")
	cfg.DB.Flags, _ = envy.MustGet("DB_FLAGS")
	cfg.DB.DSN = generateDSN(cfg.DB)

	// Mail
	cfg.Mail.Driver, _ = envy.MustGet("MAIL_DRIVER")
	cfg.Mail.Host = envy.Get("MAIL_HOST", "localhost")
	cfg.Mail.Port, _ = strconv.ParseInt(envy.Get("MAIL_PORT", "3306"), 10, 64)
	cfg.Mail.Encryption, _ = envy.MustGet("MAIL_ENCRYPTION")
	cfg.Mail.Username, _ = envy.MustGet("MAIL_USERNAME")
	cfg.Mail.Password, _ = envy.MustGet("MAIL_PASSWORD")
	cfg.Mail.FromName, _ = envy.MustGet("MAIL_NAME")
	cfg.Mail.FromAddress, _ = envy.MustGet("MAIL_ADDRESS")

	// SMS
	cfg.SMS.Driver, _ = envy.MustGet("SMS_DRIVER")
	cfg.SMS.Username, _ = envy.MustGet("SMS_USERNAME")
	cfg.SMS.Password, _ = envy.MustGet("SMS_PASSWORD")
	cfg.SMS.FromName, _ = envy.MustGet("SMS_NAME")

	// EPOINT
	cfg.Epoint.BaseURL, _ = envy.MustGet("EPOINT_BASEURL")
	cfg.Epoint.MTID, _ = strconv.ParseInt(envy.Get("EPOINT_MTID", ""), 10, 64)
	cfg.Epoint.Username, _ = envy.MustGet("EPOINT_USERNAME")
	cfg.Epoint.Password, _ = envy.MustGet("EPOINT_PASSWORD")

	// STP
	cfg.STP.BaseURL, _ = envy.MustGet("STP_API")
	cfg.STP.Username, _ = envy.MustGet("STP_USERNAME")
	cfg.STP.Password, _ = envy.MustGet("STP_PASSWORD")

	// KMS
	cfg.KMS.BaseURL, _ = envy.MustGet("KMS_API")
	cfg.KMS.ProgramID, _ = strconv.ParseInt(envy.Get("KMS_PROGRAM_ID", ""), 10, 64)
	cfg.KMS.Username, _ = envy.MustGet("KMS_USERNAME")
	cfg.KMS.Password, _ = envy.MustGet("KMS_PASSWORD")

	// Evault
	cfg.Evault.BaseURL, _ = envy.MustGet("EVAULT_API")
	cfg.Evault.ProgramID, _ = strconv.ParseInt(envy.Get("EVAULT_PROGRAM_ID", ""), 10, 64)
	cfg.Evault.PartitionID, _ = strconv.ParseInt(envy.Get("EVAULT_PARTITION_ID", ""), 10, 64)
	cfg.Evault.Username, _ = envy.MustGet("EVAULT_USERNAME")
	cfg.Evault.Password, _ = envy.MustGet("EVAULT_PASSWORD")

	// For development only
	if IsDev() {
		cfg.Mail.ToAddress = envy.Get("MAIL_ADDRESS_TO_OVERRIDE", "")
		cfg.SMS.ToMobile = envy.Get("SMS_MOBILE_TO_OVERRIDE", "")
	}

	fmt.Println("Environment: " + cfg.Environment)

	return cfg, nil
}

// Get gets the processed configuration values
func Get() Configuration {
	return cfg
}
