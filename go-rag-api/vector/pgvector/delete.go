package pgvector

import (
	"context"
)

func (s *Store) Delete(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	_, err := s.pool.Exec(ctx, "DELETE FROM documents where id = ANY($1)", ids)
	return err
}

// DeleteBySource removes all documents associated with a specific source
func (s *Store) DeleteBySource(ctx context.Context, source string) error {
	if source == "" {
		return nil
	}

	const query = `DELETE FROM documents WHERE metadata->>'source' = $1`
	_, err := s.pool.Exec(ctx, query, source)
	return err
}
