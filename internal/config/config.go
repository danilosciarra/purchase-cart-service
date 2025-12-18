package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds application configuration values
type Config struct {
	VATRate  float64
	Port     string
	WebApp   Server
	Database Database
}
type Server struct {
	HostName string
	Port     int
}
type Database struct {
	Type     string
	Host     string
	Port     int
	Username string
	Password string
	Name     string
}

// Load loads configuration from environment variables
// and applies sensible defaults for local development
func Load() *Config {
	var cfg Config
	bConfig, err := os.ReadFile("config.json")
	if err != nil {
		panic(fmt.Sprintf("Error on reading config file. Error:%s", err.Error()))
	}

	err = json.Unmarshal(bConfig, &cfg)
	if err != nil {
		panic(fmt.Sprintf("Error on unmarshalling configuration. Error:%s", err.Error()))
	}
	return &cfg
}
