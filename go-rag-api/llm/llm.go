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
	chatSDK openai.Client
	embeddingSDK openai.Client
}

func New(cfg config.Config) *Client {
	chatOpts := []option.RequestOption{}
	if cfg.ChatBaseURL != "" {
		chatOpts = append(chatOpts, option.WithBaseURL(cfg.ChatBaseURL))
	}
	if cfg.ChatAPIKey != "" {
		chatOpts = append(chatOpts, option.WithAPIKey(cfg.ChatAPIKey))
	}

	chatSDK := openai.NewClient(chatOpts...)

	embOpts := []option.RequestOption{}
    if cfg.EmbeddingBaseURL != "" {
        embOpts = append(embOpts, option.WithBaseURL(cfg.EmbeddingBaseURL))
    }
    if cfg.EmbeddingAPIKey != "" {
        embOpts = append(embOpts, option.WithAPIKey(cfg.EmbeddingAPIKey))
    }
    embeddingSDK := openai.NewClient(embOpts...)


	return &Client{
		cfg: cfg,
		chatSDK: chatSDK,
		embeddingSDK: embeddingSDK,
	}
}

func (c *Client) ChatStream(
	ctx context.Context,
	messages []Message,
	onDelta func(string),
) (Message, error) {
	stream := c.chatSDK.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Model:    c.cfg.ChatModel,
		Messages: chatMessages(messages),
	})
	defer stream.Close()

	var content strings.Builder
	role := "assistant"

	for stream.Next() {
		chunk := stream.Current()

		if len(chunk.Choices) > 0 {
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
