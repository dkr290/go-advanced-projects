package pgvector

import (
	"context"
	"fmt"

	"github.com/dkr290/go-advanced-projects/go-rag-api/vector"
	"github.com/pgvector/pgvector-go"
)

func (s *Store) Upsert(ctx context.Context, docs []vector.Document) error {
	if len(docs) == 0 {
		return nil
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

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
		metadataBytes, err := marshalMetadata(d.Metadata)
		if err != nil {
			return fmt.Errorf("metadata for %s: %w", d.ID, err)
		}
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
