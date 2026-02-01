// Package conf
package conf

import (
	"log"
	"os"
	"strings"
)

type Config struct {
	APIKey string
	Models []string
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

	models := getEnv("MODELS")
	if len(models) < 1 {
		log.Fatalln("Need models like gemini-2.5-flash in comma separate list for failback")
	}
	c.Models = strings.Split(models, ",")
}

func getEnv(key string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return ""
}
