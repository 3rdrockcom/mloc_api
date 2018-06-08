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
}

// Application contains application information
type Application struct {
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

// cfg contains the processed configuration values
var cfg Configuration

// New processes the configuration values
func New() (Configuration, error) {
	// Base
	cfg.Environiment = envy.Get("ENVIRONMENT", "development")

	// Application
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

	return cfg, nil
}

// Get gets the processed configuration values
func Get() Configuration {
	return cfg
}
