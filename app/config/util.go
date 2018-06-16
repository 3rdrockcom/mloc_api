package config

import "strconv"

// EnvDevelopment is the setting for a development environment
const EnvDevelopment string = "development"

// EnvProduction is the setting for a production environment
const EnvProduction string = "production"

// generateDSN creates a DSN from database config
func generateDSN(d Database) string {
	return d.Username + ":" + d.Password + "@tcp(" + d.Host + ":" + strconv.FormatInt(d.Port, 10) + ")/" + d.Database + "?" + d.Flags
}

// IsDev determines if the application environment is in development mode
func IsDev() bool {
	return cfg.Environiment == EnvDevelopment
}

// IsProd determines if the application environment is in production mode
func IsProd() bool {
	return cfg.Environiment == EnvProduction
}
