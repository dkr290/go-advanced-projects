package chunk

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dkr290/go-advanced-projects/go-rag-api/llm"
	"github.com/dkr290/go-advanced-projects/go-rag-api/vector"
)

// implement ingestion pipline that takes document
// from source directory and puts in the vector store back in
// 1. READ
// 2. CHUNK
// 3. EMBED
// 4. DELETE
// 5. UPSERT


const (
	defaultChunkSize    = 1000
	defaultChunkOverlap = 100
)

type Options struct {
	SourceDir    string
	ProcessedDir string
	ChunkSize    int
	ChunkOverlap int
}

// processOne reads a single file and runs the full ingestion pipeline on it.
// It delegates to preocessContent for the core logic.
func processOne(ctx context.Context,path string,opts Options,embedder llm.Embedder, store vector.Store) error {

	// Quick format check before attempting to read
	if !supportedFormat(path) {
		return fmt.Errorf("unsupported format %s", filepath.Ext(path))
	}
	// READ: load raw file contents
	raw, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %w", err)
	}
	// Delegate to the core pipeline, passing only the base filename as the source identifier
	n, err := processContent(ctx, filepath.Base(path), raw, opts, embedder, store)
	if err != nil {
		return fmt.Errorf("process %s: %w", path, err)
	}

	fmt.Printf("ingested %s: %d chunks\n", path, n)
	return nil
}

// preocessContent is the core ingestion pipeline.
// It takes the source basename and raw content, then executes CHUNK → EMBED → DELETE → UPSERT.
// ingest document from some source directoy
func processContent(
	ctx context.Context,
	source string,
	content []byte,
	opts Options,
	embedder llm.Embedder,
	store vector.Store,
) (int, error) {
	if embedder == nil {
		return 0, errors.New("embedder is required")
	}

	if store == nil {
		return 0, errors.New("vector store is required")
	}
	base := filepath.Base(source)
	if !supportedFormat(base) {
		return 0, fmt.Errorf("unsupported format %s", filepath.Ext(base))
	}

	// CHUNK: apply user-provided or default chunk size / overlap
	size := opts.ChunkSize
	if size <= 0 {
		size = defaultChunkSize
	}

	overlap := opts.ChunkOverlap
	if overlap <= 0 {
		overlap = defaultChunkOverlap
	}
	text := strings.TrimSpace(string(content))
	if text == "" {
		return 0, errors.New("file is empty")
	}
	// CHUNK the text
	// Split the text into overlapping chunks (boundary-aware: avoids cutting mid-word)
	chunks := chunk(text, size, overlap)
	if len(chunks) == 0 {
		return 0, errors.New("no chunks produced")
	}
	// EMBED: these are documents, not queries
	// so Nomic-style prefixes ("search_document: …") are applied correctly.
	vectors, err := embedder.Embed(ctx, chunks, false)
	if err != nil {
		return 0, fmt.Errorf("embed %w", err)
	}
	if len(vectors) != len(chunks) {
		return 0, fmt.Errorf("embed got %d vectors for %d chunks", len(vectors), len(chunks))
	}
	// DELETE: clear previous chunks for this source
	// DELETE: remove any prior vectors for this source so we don't leave stale data
	if err := store.DeleteBySource(ctx, base); err != nil {
		return 0, fmt.Errorf("clear previous chunks %w", err)
	}
	// UPSERT: build documents and upsert them
	// UPSERT: construct vector.Document structs and upsert them into the store.
	// Each document gets a stable ID and metadata so we can filter / debug later.
	ingestedAt := time.Now().UTC().Format(time.RFC3339)
	// build the document
	docs := make([]vector.Document, len(chunks))

	for i, c := range chunks {
		docs[i] = vector.Document{
			ID:        fmt.Sprintf("%s-chunk-%d", strings.ReplaceAll(base, ".", "_"), i),
			Content:   c,
			Metadata: map[string]string{
				"source":      base,
				"chunk_index": fmt.Sprintf("%d", i),
				"chunks": fmt.Sprintf("%d", len(chunks)),
				"ingested_at": ingestedAt,
			},
			Embedding: vectors[i],
		}
	}

	if err := store.Upsert(ctx, docs); err != nil {
		return 0, fmt.Errorf("upsert %w", err)
	}

	return len(chunks), nil
}

func supportedFormat(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".txt", ".md", "markdown":
		return true
	}
	return false
}
