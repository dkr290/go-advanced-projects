package chunk

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

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

// ingest document from some source directoy
func preocessContent(
	ctx context.Context,
	source string,
	content []byte,
	opts Options,
	embedder llm.Embedder,
	store vector.Store,
) (int, error) {
	if embedder == nil {
		return 0, errors.New("Embedder is required")
	}

	if store == nil {
		return 0, errors.New("vector store is required")
	}
	base := filepath.Base(source)
	if !supportedFormat(base) {
	return 0 , fmt.Errorf("unsupported format %s", filepath.Ext(base))
  }


	size := opts.ChunkSize
	if size <=0 {

		size = defaultChunkSize
	}

	overlap := opts.ChunkOverlap
	if overlap <= 0 {
     overlap = defaultChunkOverlap
	}

}

func supportedFormat(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".txt", ".md", "markdown":
		return true
	}
	return false
}
