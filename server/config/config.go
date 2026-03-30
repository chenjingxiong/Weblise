package config

import (
	"fmt"
	"os"
)

type Config struct {
	HTTPPort string
	WSPort   string
	DBPath   string
}

func Load() *Config {
	return &Config{
		HTTPPort: getEnv("HTTP_PORT", "8080"),
		WSPort:   getEnv("WS_PORT", "8443"),
		DBPath:   getEnv("DB_PATH", "./data/weblise.db"),
	}
}

func (c *Config) HTTPAddr() string {
	return fmt.Sprintf(":%s", c.HTTPPort)
}

func (c *Config) WSAddr() string {
	return fmt.Sprintf(":%s", c.WSPort)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
