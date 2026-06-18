package ai

import (
	"context"
	"time"

	"github.com/agenticgokit/agenticgokit/v1beta"
	"github.com/example/go-code-agent/pkg/tools"

	_ "github.com/agenticgokit/agenticgokit/plugins/llm/openai" // registers the "openai" provider
)

const basePrompt = "You are an autonomous coding assistant. You use tools to read, search, and write files in the local running workspace."

func SetupAgent(modelName string, AIUrl string) (v1beta.Agent, error) {
	ts := tools.NewToolset(basePrompt).
		Add("search_workspace_files", "Lists every file path in the workspace. No arguments.", nil,
			func(ctx context.Context, _ map[string]interface{}) (string, error) {
				return tools.SearchWorkspaceFiles(ctx)
			}).
		Add("read_file_content", "Reads the full text of a file.",
			map[string]interface{}{"path": "string - file to read"},
			func(ctx context.Context, a map[string]interface{}) (string, error) {
				return tools.ReadFileContent(ctx, tools.Arg(a, "path"))
			}).
		Add("write_file_content", "Creates folders if needed and writes text to a file.",
			map[string]interface{}{"path": "string - destination", "content": "string - file body"},
			func(ctx context.Context, a map[string]interface{}) (string, error) {
				return tools.WriteFileContent(ctx, tools.Arg(a, "path"), tools.Arg(a, "content"))
			})

	return v1beta.NewBuilder("CodeAgent").
		WithConfig(&v1beta.Config{
			Name:         "CodeAgent",
			SystemPrompt: basePrompt,
			LLM:          v1beta.LLMConfig{Provider: "openai", Model: modelName, BaseURL: AIUrl},
			Tools:        &v1beta.ToolsConfig{Enabled: true, MaxRetries: 3, Timeout: 30 * time.Second},
		}).
		WithHandler(ts.Handler(5)).
		Build()
}

