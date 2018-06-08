package config

import "fmt"

// Version is the application semantic version
var Version string

// Build is the application build version
var Build string

func init() {
	fmt.Println("Version: ", Version)
	fmt.Println("Build: ", Build)
}
