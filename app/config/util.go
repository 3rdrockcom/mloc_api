package config

import "strconv"

// generateDSN creates a DSN from database config
func generateDSN(d Database) string {
	return d.Username + ":" + d.Password + "@tcp(" + d.Host + ":" + strconv.FormatInt(d.Port, 10) + ")/" + d.Database + "?" + d.Flags
}
