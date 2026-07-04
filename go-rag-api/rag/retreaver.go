package rag

import (
	"context"
	"fmt"

	"github.com/dkr290/go-advanced-projects/go-rag-api/llm"
	"github.com/dkr290/go-advanced-projects/go-rag-api/vector"
)

const defaultTopK = 5 // number of chunks from vector store

type Options struct {
	TopK     int
	Rewriter *Rewriter
}

type Retriever struct {
	embedder llm.Embedder
	store    vector.Store
	rewriter *Rewriter
	topK     int
}

func New(embedder llm.Embedder, store vector.Store, opts Options) *Retriever {
	topK := opts.TopK
	if topK <= 0 {
		topK = defaultTopK
	}

	return &Retriever{
		embedder: embedder,
		store:    store,
		rewriter: opts.Rewriter,
		topK:     topK,
	}
}


func getLastUserMessage(history []llm.Message) string {
	for i := len(history) - 1; i >= 0; i-- {
		if history[i].Role == "user" {
			return history[i].Content
		}
	}
	return ""
}

// Retrieve queries the vector store using the user's intent and returns formatted context.
func (r *Retriever) Retrieve(ctx context.Context, history []llm.Message) (string, error) {
	// Build the query
	var query string
	if r.rewriter != nil {
		var err error
		query, err = r.rewriter.Rewrite(ctx, history)
		if err != nil {
			return "", fmt.Errorf("rewrite query: %w", err)
		}
	} else {
		query = getLastUserMessage(history)
	}
	if query == "" {
		return "", nil
	}

	// Embed the query
	embeddings, err := r.embedder.Embed(ctx, []string{query}, true)
	if err != nil {
		return "", fmt.Errorf("embed query: %w", err)
	}
	if len(embeddings) == 0 {
		return "", nil
	}

	// Search the vector store
	hits, err := r.store.Query(ctx, embeddings[0], r.topK)
	if err != nil {
		return "", fmt.Errorf("vector store query: %w", err)
	}
	

	return formatContext(hits), nil
}
