package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	BaseURL          string
	APIKey           string
	Model            string
	EmbeddingModel   string 
	SystemPromptFile string
	DatabaseURL      string
	EmbeddingDIM     int
}

func Load() Config {
	_ = godotenv.Load()

	cfg := Config{
		BaseURL:          os.Getenv("BASE_URL"),
		APIKey:           os.Getenv("API_KEY"),
		Model:            os.Getenv("MODEL"),
		EmbeddingModel:   os.Getenv("EMBEDDING_MODEL"), 
		SystemPromptFile: os.Getenv("SYSTEM_PROMPT_FILE"),
		DatabaseURL:      os.Getenv("DATABASE_URL"),
		EmbeddingDIM:     atoiOr(os.Getenv("EMBEDDING_DIM"), 0),
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = "http://localai:8080"
	}
	if cfg.Model == "" {
		cfg.Model = "supergemma4-26b-uncensored-v2"
	}

	if cfg.EmbeddingModel == "" { 
		cfg.EmbeddingModel = "nomic-embed-text-v1.5"
	}
	if cfg.SystemPromptFile == "" {
		cfg.SystemPromptFile = "./prompts/system-custom.md"
	}
	if cfg.DatabaseURL == "" {
		cfg.DatabaseURL = "postgresql://rag:rag@localhost:5432/rag?sslmode=disable"
	}
	if cfg.EmbeddingDIM == 0 {

		cfg.EmbeddingDIM = 768
	}

	return cfg
}

func atoiOr(s string, fallback int) int {
	if s == "" {
		return fallback
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}
	return n
}
