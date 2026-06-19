package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// chat configuration
	ChatBaseURL string
	ChatAPIKey  string
	ChatModel   string

	// Embedding configuration (falls back to Chat if empty)
	EmbeddingBaseURL string
	EmbeddingAPIKey  string
	EmbeddingModel   string

	SystemPromptFile string
	DatabaseURL      string
	EmbeddingDIM     int
	IngestDir        string
	ProcessedDir     string
}

func Load() Config {
	_ = godotenv.Load()

	chatBaseURL := envOrDefault("CHAT_BASE_URL", "BASE_URL", "")
	chatAPIKey := envOrDefault("CHAT_API_KEY", "API_KEY", "")
	chatModel := envOrDefault("CHAT_MODEL", "MODEL", "")
	embBaseURL := envOrDefault("EMBEDDING_BASE_URL", "", "")
	embAPIKey := envOrDefault("EMBEDDING_API_KEY", "", "")
	embModel := envOrDefault("EMBEDDING_MODEL", "", "")

	cfg := Config{
		ChatBaseURL:      chatBaseURL,
		ChatAPIKey:       chatAPIKey,
		ChatModel:        chatModel,
		EmbeddingBaseURL: embBaseURL,
		EmbeddingAPIKey:  embAPIKey,
		EmbeddingModel:   embModel,
		SystemPromptFile: os.Getenv("SYSTEM_PROMPT_FILE"),
		DatabaseURL:      os.Getenv("DATABASE_URL"),
		EmbeddingDIM:     atoiOr(os.Getenv("EMBEDDING_DIM"), 0),
		IngestDir:        os.Getenv("INGEST_DIR"),
		ProcessedDir:     os.Getenv("PROCESSED_DIR"),
	}

	// Chat defaults
	if cfg.ChatBaseURL == "" {
		cfg.ChatBaseURL = "http://localai:8080"
	}
	if cfg.ChatModel == "" {
		cfg.ChatModel = "supergemma4-26b-uncensored-v2"
	}
	if cfg.ChatAPIKey == "" {
		cfg.ChatAPIKey = ""
	}

	// Embedding defaults
	if cfg.EmbeddingModel == "" {
		cfg.EmbeddingModel = "nomic-embed-text-v1.5"
	}

	// Fallback: Embedding-specific values default to Chat values
	if cfg.EmbeddingBaseURL == "" {
		cfg.EmbeddingBaseURL = cfg.ChatBaseURL
	}
	if cfg.EmbeddingAPIKey == "" {
		cfg.EmbeddingAPIKey = cfg.ChatAPIKey
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
	if cfg.IngestDir == "" {
		cfg.IngestDir = "./documents"
	}
	if cfg.ProcessedDir == "" {
		cfg.ProcessedDir = "./documents/processed"
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

func envOrDefault(first, fallback, defaultVal string) string {
	if v := os.Getenv(first); v != "" {
		return v
	}
	if v := os.Getenv(fallback); v != "" {
		return v
	}
	return defaultVal
}
