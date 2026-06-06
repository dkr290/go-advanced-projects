package vector

import "context"

type Document struct {
	ID        string
	Content   string            // just to see the content
	Metadata  map[string]string // structured data metadata
	Embedding []float32
}

// Result is  one hit for similarity search
type Result struct {
	Document
	Score float32
}

// Store interface to be implemented easy to switch
type Store interface {
	Upsert(ctx context.Context, docs []Document) error
}
