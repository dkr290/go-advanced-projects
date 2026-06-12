package pgvector

import (
	"context"
	"fmt"

	"github.com/dkr290/go-advanced-projects/go-rag-api/vector"
	"github.com/pgvector/pgvector-go"
)

// Upsert inserts new documents or updates existing ones in the database.
// It operates within a single transaction to ensure atomicity — either all
// documents are upserted successfully or none are.
func (s *Store) Upsert(ctx context.Context, docs []vector.Document) error {
	// Early return if there is nothing to process.
	if len(docs) == 0 {
		return nil
	}

	// Start a transaction.
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Rollback on exit if the transaction was not already committed.
	// This is safe because Rollback on an already-committed transaction is a no-op.
	// this here 
	defer tx.Rollback(ctx)

	// Upsert SQL: inserts a new row or updates the matching row on ID conflict.
	// ON CONFLICT (id) DO UPDATE sets the three mutable fields and refreshes
	// created_at to the current time.
	const upsertSQL = `
		INSERT INTO documents (id, content, metadata, embedding)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			content = EXCLUDED.content,
			metadata = EXCLUDED.metadata,
			embedding = EXCLUDED.embedding,
			created_at = NOW()
	`
	for _, d := range docs {

		// Marshal the map-based metadata into a JSONB-compatible byte slice.
		metadataBytes, err := marshalMetadata(d.Metadata)
		if err != nil {
			return fmt.Errorf("metadata for %s: %w", d.ID, err)
		}

		// Execute the upsert for this document.
		_, err = tx.Exec(
			ctx,
			upsertSQL,
			d.ID,
			d.Content,
			metadataBytes,
			pgvector.NewVector(d.Embedding),
		)
		if err != nil {
			return fmt.Errorf("failed to upsert document %s: %w", d.ID, err)
		}

	}

	return tx.Commit(ctx)
}
