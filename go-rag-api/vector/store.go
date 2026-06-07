// Package vector for vector embedding operations
package vector

import "context"

// Document represents a document with its vector embedding
type Document struct {
	// ID uniquely identifies the document
	ID string
	// Content contains the raw text content of the document
	Content string
	// Metadata contains structured data associated with the document
	Metadata map[string]string
	// Embedding is the vector representation of the document
	Embedding []float32
}

// Result represents a search result containing a document and its similarity score
type Result struct {
	// Document is the retrieved document
	Document
	// Score represents the similarity score between the query and document
	Score float32
}

// Store interface defines the contract for vector storage operations
// This interface allows easy switching between different vector database implementations
// such as PostgreSQL with pgvector, Weaviate, or other vector databases
type Store interface {
	// Upsert inserts or updates documents in the vector store
	// Documents with the same ID will be updated with new content and embeddings
	Upsert(ctx context.Context, docs []Document) error

	// Query performs a similarity search using vector embeddings
	// Returns up to topK most similar documents sorted by similarity score
	// The embedding parameter should contain the query vector
	// Returns an error if the operation fails
	Query(ctx context.Context, embedding []float32, topK int) ([]Result, error)

	// Delete removes documents with the specified IDs from the vector store
	Delete(ctx context.Context, ids []string) error

	// DeleteBySource removes all documents associated with a specific source
	DeleteBySource(ctx context.Context, source string) error

	// Close releases any resources held by the store
	Close() error
}
