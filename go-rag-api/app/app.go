package app

import (
	"context"
	"log"
	"os"

	"github.com/dkr290/go-advanced-projects/go-rag-api/chat"
	"github.com/dkr290/go-advanced-projects/go-rag-api/config"
	"github.com/dkr290/go-advanced-projects/go-rag-api/llm"
	"github.com/dkr290/go-advanced-projects/go-rag-api/vector"
	"github.com/dkr290/go-advanced-projects/go-rag-api/vector/pgvector"
)

func Run(ctx context.Context, cfg config.Config) error {
	// define some better logging
	logger := log.New(os.Stderr, "[rag] ", log.LstdFlags)

	client := llm.New(cfg)
	store, err := openStore(ctx, cfg)
	if err != nil {
		logger.Printf("vector store disabled: %v", err)
	}

	if store != nil {
     defer store.Close()
		 logger.Println("vector store ready")

	}
	return chat.RunREPL(ctx, client, chat.Options{
		SystemPromptFile: cfg.SystemPromptFile,
	})
}

func openStore(ctx context.Context, cfg config.Config) (vector.Store, error) {
	s, err := pgvector.New(ctx, pgvector.Options{
		DSN:          cfg.DatabaseURL,
		EmbeddingDim: cfg.EmbeddingDIM,
	})
	if err != nil {
		return nil, err
	}

	return s, nil
}
