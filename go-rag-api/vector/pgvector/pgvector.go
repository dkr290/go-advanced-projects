// Package pgvector is the package for implementing
package pgvector

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxvec "github.com/pgvector/pgvector-go/pgx"
)

// types specific to this postgres store

type Options struct {
	DSN          string
	EmbeddingDim int
}

type Store struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, opts Options) (*Store, error) {
	if opts.DSN == "" {
		return nil, fmt.Errorf("DSN is required")
	}
	if opts.EmbeddingDim <= 0 {
		return nil, errors.New("pgvector: EmbeddingDim must be > 0")
	}

	var pool *pgxpool.Pool
	var err error

	// retry the connection up to 10 times
	maxRetries := 10
	delay := 1 * time.Second

	poolConfig, err := pgxpool.ParseConfig(opts.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}

	poolConfig.AfterConnect = func(ctx context.Context, c *pgx.Conn) error {
		return pgxvec.RegisterTypes(ctx, c)
	}

	for i := range maxRetries {

		pool, err = pgxpool.NewWithConfig(ctx, poolConfig)
		if err != nil {
			if i == maxRetries-1 {
				return nil, fmt.Errorf(
					"failed to create connection pool after %d attempts: %w",
					maxRetries,
					err,
				)
			}

			time.Sleep(delay)
			continue
		}

		// Test the connection
		if err = pool.Ping(ctx); err != nil {
			pool.Close()
			if i == maxRetries-1 {
				return nil, fmt.Errorf(
					"failed to ping database after %d attempts: %w",
					maxRetries,
					err,
				)
			}
			time.Sleep(delay)
			continue
		}
		// Check if pgvector extension is installed
		if err = checkPgVectorExtension(ctx, pool); err != nil {
			pool.Close()
			return nil, fmt.Errorf("pgvector extension check failed: %w", err)
		}
		break
	}

	s := &Store{pool: pool}
	if err = s.migrate(ctx, opts.EmbeddingDim); err != nil {
		pool.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return s, nil
}

// checkPgVectorExtension verifies that the pgvector extension is installed
func checkPgVectorExtension(ctx context.Context, pool *pgxpool.Pool) error {
	const checkExtensionSQL = `
		SELECT
			extname
		FROM
			pg_extension
		WHERE
			extname = 'vector'
	`

	var extensionName string
	err := pool.QueryRow(ctx, checkExtensionSQL).Scan(&extensionName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf(
				"pgvector extension is not installed. Please run 'CREATE EXTENSION IF NOT EXISTS vector;' in your database",
			)
		}
		return fmt.Errorf("failed to check pgvector extension: %w", err)
	}

	return nil
}

func (s *Store) migrate(ctx context.Context, embeddingDim int) error {
	const migrateSQL = `
		CREATE TABLE IF NOT EXISTS documents (
			id TEXT PRIMARY KEY,
			content TEXT NOT NULL,
			metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
			embedding VECTOR(%d) NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS idx_documents_embedding ON documents USING hnsw (embedding vector_cosine_ops);
CREATE INDEX IF NOT EXISTS idx_documents_metadata ON documents USING GIN (metadata);

	`
	_, err := s.pool.Exec(ctx, fmt.Sprintf(migrateSQL, embeddingDim))
	return err
}

func marshalMetadata(m map[string]string) ([]byte, error) {
	if len(m) == 0 {
		return []byte("{}"), nil
	}

	return json.Marshal(m)
}

// Close releases any resources held by the store
func (s *Store) Close() error {
	s.pool.Close()
	return nil
}
