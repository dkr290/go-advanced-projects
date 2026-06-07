package pgvector

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dkr290/go-advanced-projects/go-rag-api/vector"
	"github.com/pgvector/pgvector-go"
)

func (s *Store) Query(ctx context.Context, embedding []float32, topK int) ([]vector.Result, error) {
	if topK <= 0 {
		return nil, nil
	}

	const querySQL = `
		SELECT 
			id,
			content,
			metadata,
			embedding <=> $1) AS distance
		FROM documents 
		ORDER BY embedding <=> $1
		LIMIT $2
	`


	rows, err := s.pool.Query(ctx, querySQL, pgvector.NewVector(embedding), topK)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var results []vector.Result
	for rows.Next(){
    var r vector.Result
		var metaRaw []byte
		var distance float64

		if err := rows.Scan(&r.ID,&r.Content,&metaRaw,&distance);err != nil {


			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if err := unMarshalMetadata(metaRaw, &r.Metadata);err != nil {

			return nil , fmt.Errorf("metadata for %s: %w", r.ID,err)
		}

		r.Score = float32(1 - distance)
    results = append(results, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration failed: %w", err)
	}
	return  results,nil
}


func unMarshalMetadata(raw []byte,dst *map[string]string) error {
   if len(raw) == 0 {
     *dst = nil 
		 return nil
	 }
 return json.Unmarshal(raw, dst)
}
