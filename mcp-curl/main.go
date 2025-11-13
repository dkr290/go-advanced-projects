package main

import (
	"context"
	"log"
	"log/slog"
	"os/exec"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Input struct {
	URL string `json:"url" description:"url of the webpage to fetch"`
}

type Output struct {
	Content string `json:"content" description:"fetched webpage content"`
}

func main() {
	s := mcp.NewServer(
		&mcp.Implementation{Name: "mcp-curl", Version: "1.0.0"},
		&mcp.ServerOptions{Logger: slog.Default()},
	)

	toolOptions := mcp.Tool{
		Description: "fetch a webpage using curl",
		Name:        "Curl tool",
	}
	mcp.AddTool(s, &toolOptions, curlContent)
	if err := s.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}

func curlContent(ctx context.Context, _ *mcp.CallToolRequest, input Input) (
	*mcp.CallToolResult,
	Output,
	error,
) {
	url := input.URL

	cmd := exec.Command("curl", "-s", url)
	output, err := cmd.Output()
	if err != nil {
		return &mcp.CallToolResult{}, Output{}, err
	}

	content := string(output)
	return nil, Output{Content: content}, nil
}
