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
	Debug  string
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
	c.Models = getModels()
	if len(c.Models) == 0 {
		log.Fatalln("MODELS is missing and needs to be set as env")
	}

	if debugFlag := getEnv("DEBUG"); debugFlag != "" {
		c.Debug = debugFlag
	} else {
		c.Debug = "false"
	}
}

func getEnv(key string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return ""
}

func getModels() []string {
	val := os.Getenv("MODELS")
	if val == "" {
		return []string{"gemini-2.0-flash"}
	}
	models := strings.Split(val, ",")
	for i := range models {
		models[i] = strings.TrimSpace(models[i])
	}
	return models
}
