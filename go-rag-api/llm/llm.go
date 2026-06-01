// Package llm
package llm

import (
	"context"
	"strings"

	"github.com/dkr290/go-advanced-projects/go-rag-api/config"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Client struct {
	cfg config.Config
	sdk openai.Client
}

func New(cfg config.Config) *Client {
	opts := []option.RequestOption{}
	if cfg.BaseURL != "" {
		opts = append(opts, option.WithBaseURL(cfg.BaseURL))
	}
	if cfg.APIKey != "" {
		opts = append(opts, option.WithAPIKey(cfg.APIKey))
	}

	sdk := openai.NewClient(opts...)

	return &Client{
		cfg: cfg,
		sdk: sdk,
	}
}

func (c *Client) ChatStream(
	ctx context.Context,
	messages []Message,
	onDelta func(string),
) (Message, error) {
	stream := c.sdk.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Model:    c.cfg.Model,
		Messages: chatMessages(messages),
	})
	defer stream.Close()

	var content strings.Builder
	role := "assistant"

	for stream.Next() {
		chunk := stream.Current()

		if chunk.Choices != nil && len(chunk.Choices) > 0 {
			delta := chunk.Choices[0].Delta
			if delta.Role != "" {
				role = delta.Role
			}
			if delta.Content != "" {
				content.WriteString(delta.Content)
				if onDelta != nil {
					onDelta(delta.Content)
				}
			}

		}
	}

	if err := stream.Err(); err != nil {
		return Message{}, err
	}

	return Message{Role: role, Content: content.String()}, nil
}

func chatMessages(messages []Message) []openai.ChatCompletionMessageParamUnion {
	out := make([]openai.ChatCompletionMessageParamUnion, 0, len(messages))

	for _, msg := range messages {
		switch msg.Role {
		case "system":
			out = append(out, openai.SystemMessage(msg.Content))
		case "assistant":
			out = append(out, openai.AssistantMessage(msg.Content))
    default:
       out = append(out, openai.UserMessage(msg.Content))
		}
	}
	return out
}
