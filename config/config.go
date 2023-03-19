package config

import (
	"os"
)

// Config is exported via godotenv
func Config(key string) string {
	return os.Getenv(key)
}
