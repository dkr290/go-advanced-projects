// Package handlers exposes thin MCP tool handlers that delegate business
// logic to pkg/utils.
package handlers

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"duckduckgo-mcp/pkg/utils"
)

// SearchInput is the input schema for the "duckduckgo_search" tool.
type SearchInput struct {
	Query      string `json:"query" jsonschema:"the search query / topic to look up on the internet"`
	MaxResults int    `json:"max_results,omitempty" jsonschema:"maximum number of results to return (default 10, max 25)"`
}

// SearchOutput is the output schema for the "duckduckgo_search" tool.
type SearchOutput struct {
	Query   string               `json:"query"`
	Results []utils.SearchResult `json:"results"`
	Count   int                  `json:"count"`
}

// SearchHandler wires a DuckDuckGoClient into an MCP tool handler function.
type SearchHandler struct {
	client            *utils.DuckDuckGoClient
	defaultMaxResults int
}

// NewSearchHandler builds a SearchHandler using the given client and the
// default max results to apply when the caller does not specify one.
func NewSearchHandler(client *utils.DuckDuckGoClient, defaultMaxResults int) *SearchHandler {
	return &SearchHandler{client: client, defaultMaxResults: defaultMaxResults}
}

const maxAllowedResults = 25

// Handle implements the MCP tool function signature expected by
// mcp.AddTool: it performs a DuckDuckGo search and returns structured
// results back to the calling model.
func (h *SearchHandler) Handle(ctx context.Context, _ *mcp.CallToolRequest, input SearchInput) (
	*mcp.CallToolResult,
	SearchOutput,
	error,
) {
	limit := input.MaxResults
	if limit <= 0 {
		limit = h.defaultMaxResults
	}
	if limit > maxAllowedResults {
		limit = maxAllowedResults
	}

	results, err := h.client.Search(ctx, input.Query, limit)
	if err != nil {
		return nil, SearchOutput{}, fmt.Errorf("duckduckgo search failed: %w", err)
	}

	return nil, SearchOutput{
		Query:   input.Query,
		Results: results,
		Count:   len(results),
	}, nil
}
