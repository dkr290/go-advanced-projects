package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BaseURL          string
	APIKey           string
	Model            string
	SystemPromptFile string
}

func Load() Config {
	_ = godotenv.Load()

	cfg := Config{
		BaseURL:          os.Getenv("BASE_URL"),
		APIKey:           os.Getenv("API_KEY"),
		Model:            os.Getenv("MODEL"),
		SystemPromptFile: os.Getenv("SYSTEM_PROMPT_FILE"),
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://localai:8080"
	}
	if cfg.Model == "" {
		cfg.Model = "supergemma4-26b-uncensored-v2"
	}

	return cfg
}
