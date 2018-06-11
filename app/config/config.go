package config

import (
	"strconv"

	"github.com/gobuffalo/envy"
)

// Configuration contains the application configuration
type Configuration struct {
	Environiment string
	Application  Application
	Server       Server
	DB           Database
	Mail         Mail
	SMS          SMS
}

// Application contains application information
type Application struct {
	Name    string
	Build   string
	Version string
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

// cfg contains the processed configuration values
var cfg Configuration

// New processes the configuration values
func New() (Configuration, error) {
	// Base
	cfg.Environiment = envy.Get("ENVIRONMENT", "development")

	// Application
	cfg.Application.Name = envy.Get("NAME", "app")
	cfg.Application.Build = Build
	cfg.Application.Version = Version

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

	// For development only
	if IsDev() {
		cfg.Mail.ToAddress = envy.Get("MAIL_ADDRESS_TO_OVERRIDE", "")
		cfg.SMS.ToMobile = envy.Get("SMS_MOBILE_TO_OVERRIDE", "")
	}

	return cfg, nil
}

// Get gets the processed configuration values
func Get() Configuration {
	return cfg
}
