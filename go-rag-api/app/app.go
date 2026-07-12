package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/dkr290/go-advanced-projects/go-rag-api/chat"
	"github.com/dkr290/go-advanced-projects/go-rag-api/chunk"
	"github.com/dkr290/go-advanced-projects/go-rag-api/config"
	"github.com/dkr290/go-advanced-projects/go-rag-api/llm"
	"github.com/dkr290/go-advanced-projects/go-rag-api/rag"
	"github.com/dkr290/go-advanced-projects/go-rag-api/vector"
	"github.com/dkr290/go-advanced-projects/go-rag-api/vector/pgvector"
	"github.com/dkr290/go-advanced-projects/go-rag-api/web"
)

func Run(parentCtx context.Context, cfg config.Config) error {
	// define some better logging
	logger := log.New(os.Stderr, "[rag] ", log.LstdFlags)
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	client := llm.New(cfg)

	store, err := openStore(ctx, cfg)
	if err != nil {
		return fmt.Errorf("vector store disabled: %v", err)
	}
	defer store.Close()

	logger.Println("vector store ready")

	logger.Printf("chat model=%q base_url=%q", cfg.ChatModel, cfg.ChatBaseURL)
	logger.Printf("embedding model=%q base_url=%q", cfg.EmbeddingModel, cfg.EmbeddingBaseURL)

	var wg sync.WaitGroup
	if store != nil {
		wg.Add(1)
		go func() {
			opts := chunk.Options{
				SourceDir:    cfg.IngestDir,
				ProcessedDir: cfg.ProcessedDir,
			}
			if err := chunk.Watch(
				ctx,
				opts,
				client,
				store,
				logger,
			); err != nil &&
				ctx.Err() == nil {
				logger.Printf("watcher stopped: %v", err)
			}

			defer wg.Done()
		}()
		logger.Printf("watching %s for new documents", cfg.IngestDir)
	}
	// get a retreiver , which means we need a rewriter.

	var retriever *rag.Retriever
	if store != nil {
		retriever = rag.New(client, store, rag.Options{
			TopK:     5,
			Rewriter: rag.NewRewriter(client),
		})
	}

	srv, err := web.New(client, retriever, web.Options{
		Addr:             cfg.HttpAddr,
		SystemPromptFile: cfg.SystemPromptFile,
		Store:            store,
		ProcessDir:       cfg.ProcessedDir,
		ImagesDir:        cfg.ImageDir,
	})
	if err != nil {
		return fmt.Errorf("web.New %w", err)
	}
	app := srv.Routes()

	lis, err := net.Listen("tcp", cfg.HttpAddr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Listener(lis)
	}()

	replErr := chat.RunREPL(ctx, client, retriever, chat.Options{
		SystemPromptFile: cfg.SystemPromptFile,
	})
	cancel()
	wg.Wait()
	return replErr
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
