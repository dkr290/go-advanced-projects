// Package conf
package conf

import (
	"log"
	"os"
)

type Config struct {
	APIKey string
}

func LoadConfig() *Config {
	c := &Config{}
	c.GetFlags()
	return c
}

func (c *Config) GetFlags() {
	if apiKey := getEnv("API_KEY"); apiKey != "" {
		c.APIKey = apiKey
	}
	if c.APIKey == "" {
		log.Fatalln("API_KEY is missing and needs to be set as env")
	}
}

func getEnv(key string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return ""
}
