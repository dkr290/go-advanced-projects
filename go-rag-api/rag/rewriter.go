package rag

import (
	"context"
	"strings"

	"github.com/dkr290/go-advanced-projects/go-rag-api/llm"
)

const rewriteSystemPrompt = `You rewrite the user's latest message into a standalone search query.

Given the conversation, output a single search query that:
- Captures the topic and intent of the latest user message.
- Resolves pronouns and references using prior turns ("it", "they", "that one").
- Stays concise - keywords and short phrases, not full sentences.

If the latest user message already stands on its own with no references to prior turns, output it verbatim.

Output only the query. No preamble, no quotes, no explanation.`

type Rewriter struct {
	client *llm.Client
}

func NewRewriter(client *llm.Client) *Rewriter {
	return &Rewriter{client: client}
}

func (r *Rewriter) Rewrite(ctx context.Context, history []llm.Message) (string, error) {
	last := getLastUserMessage(history)
	if last == "" {
		return "", nil
	}

	if len(history) <= 1 {
		return last, nil
	}

	messages := []llm.Message{
		{Role: "system", Content: rewriteSystemPrompt},
		{Role: "user", Content: formatConversation(history)},
	}

	resp, err := r.client.ChatStream(ctx, messages, nil)
	if err != nil {
		return "", err
	}
	return resp.Content, nil
}


func formatConversation(history []llm.Message) string {
	var sb strings.Builder
	sb.WriteString("Conversation so far:\n\n")
	for _, m := range history {
		switch m.Role {
		case "user":
			sb.WriteString("User: ")
		case "assistant":
			sb.WriteString("Assistant: ")
		default:
			continue
		}
		sb.WriteString(m.Content)
		sb.WriteString("\n\n")
	}
	return sb.String()
}
