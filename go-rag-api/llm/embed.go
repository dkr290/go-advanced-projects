package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/openai/openai-go/v3"
)

type Embedder interface {
	Embed(ctx context.Context, text []string, isQuery bool) ([][]float32, error)
}

func (c *Client) Embed(ctx context.Context, texts []string, isQuery bool) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, nil
	}

	prefixedTexts := make([]string, len(texts))

	// Apply Nomic prefixes only if using a Nomic model
	isNomic := strings.Contains(strings.ToLower(c.cfg.EmbeddingModel), "nomic")

	for i, t := range texts {
		if isNomic && isQuery {
			prefixedTexts[i] = "Search Query: " + t
		} else if isNomic && !isQuery {
			prefixedTexts[i] = "Search Article: " + t
		} else {
			prefixedTexts[i] = t
		}
	}

	resp, err := c.sdk.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Model: c.cfg.EmbeddingModel,
		Input: openai.EmbeddingNewParamsInputUnion{OfArrayOfStrings: texts},
	})
	if err != nil {
		return nil, err
	}
	if len(resp.Data) != len(texts) {
		return nil, fmt.Errorf("expected %d embeddings, got %d", len(texts), len(resp.Data))
	}
	out := make([][]float32, len(texts))
	for i, emb := range resp.Data {
		f32emb := make([]float32, len(emb.Embedding))
		for j, v := range emb.Embedding {
			f32emb[j] = float32(v)
		}
		out[i] = f32emb
	}
	return out, nil
}
