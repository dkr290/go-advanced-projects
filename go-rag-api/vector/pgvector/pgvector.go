// Package pgvector
package pgvector

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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

	pool, err := pgxpool.New(ctx, opts.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
// Check if pgvector extension is installed
	if err := checkPgVectorExtension(ctx, pool); err != nil {
		return nil, fmt.Errorf("pgvector extension check failed: %w", err)
	}

	return &Store{
		pool: pool,
	}, nil
}
// checkPgVectorExtension verifies that the pgvector extension is installed
func checkPgVectorExtension(ctx context.Context, pool *pgxpool.Pool) error {
	const checkExtensionSQL = `
		SELECT 
			extension_name 
		FROM 
			pg_extension 
		WHERE 
			extension_name = 'vector'
	`

	var extensionName string
	err := pool.QueryRow(ctx, checkExtensionSQL).Scan(&extensionName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("pgvector extension is not installed. Please run 'CREATE EXTENSION IF NOT EXISTS vector;' in your database")
		}
		return fmt.Errorf("failed to check pgvector extension: %w", err)
	}

	return nil
}

