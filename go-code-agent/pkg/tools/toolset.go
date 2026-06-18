// Package tools for varios helper tools
package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/agenticgokit/agenticgokit/v1beta"
)

// Func is the generic signature every registered tool is adapted to.
type Func func(ctx context.Context, args map[string]interface{}) (string, error)

type entry struct {
	info v1beta.ToolInfo
	fn   Func
}

// Toolset collects tools and builds a v1beta handler that runs the call loop.
type Toolset struct {
	system  string
	entries map[string]entry
	order   []string
}

func NewToolset(systemPrompt string) *Toolset {
	return &Toolset{system: systemPrompt, entries: map[string]entry{}}
}

// Add registers one tool. Returns the set for chaining.
func (t *Toolset) Add(name, desc string, params map[string]interface{}, fn Func) *Toolset {
	t.entries[name] = entry{info: v1beta.ToolInfo{Name: name, Description: desc, Parameters: params}, fn: fn}
	t.order = append(t.order, name)
	return t
}

// Handler returns a v1beta.HandlerFunc that lets the LLM call the registered tools.
func (t *Toolset) Handler(maxIterations int) func(context.Context, string, *v1beta.Capabilities) (string, error) {
	infos := make([]v1beta.ToolInfo, 0, len(t.order))
	for _, n := range t.order {
		infos = append(infos, t.entries[n].info)
	}
	prompt := t.system + v1beta.FormatToolsPromptForLLM(infos)

	return func(ctx context.Context, input string, caps *v1beta.Capabilities) (string, error) {
		conversation := input
		for i := 0; i < maxIterations; i++ {
			reply, err := caps.LLM(prompt, conversation)
			if err != nil {
				return "", fmt.Errorf("llm call failed: %w", err)
			}
			calls := v1beta.ParseLLMToolCalls(reply)
			if len(calls) == 0 {
				return reply, nil // final answer
			}

			var sb strings.Builder
			sb.WriteString(conversation)
			sb.WriteString("\n\n[tool results]\n")
			for _, call := range calls {
				name, _ := call["name"].(string)
				args, _ := call["args"].(map[string]interface{})
				e, ok := t.entries[name]
				if !ok {
					fmt.Fprintf(&sb, "- %s: unknown tool\n", name)
					continue
				}
				out, err := e.fn(ctx, args)
				if err != nil {
					fmt.Fprintf(&sb, "- %s error: %v\n", name, err)
					continue
				}
				fmt.Fprintf(&sb, "- %s:\n%s\n", name, out)
			}
			conversation = sb.String()
		}
		return caps.LLM(prompt, conversation+"\n\nProvide your final answer now without requesting more tools.")
	}
}

// Arg is a small helper to read a string argument.
func Arg(args map[string]interface{}, key string) string {
	if v, ok := args[key]; ok {
		return fmt.Sprintf("%v", v)
	}
	return ""
}
// SearchWorkspaceFiles walks the current working directory and returns all relative file paths.
func SearchWorkspaceFiles(_ context.Context) (string, error) {
	root, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getwd: %w", err)
	}
	var paths []string
	err = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && strings.HasPrefix(d.Name(), ".") {
			return filepath.SkipDir
		}
		if !d.IsDir() {
			rel, _ := filepath.Rel(root, path)
			paths = append(paths, rel)
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("walk: %w", err)
	}
	return strings.Join(paths, "\n"), nil
}

// ReadFileContent reads and returns the full contents of a file.
func ReadFileContent(_ context.Context, path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path is required")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", path, err)
	}
	return string(data), nil
}

// WriteFileContent creates parent directories as needed and writes content to path.
func WriteFileContent(_ context.Context, path, content string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path is required")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", fmt.Errorf("mkdir: %w", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return "", fmt.Errorf("write %s: %w", path, err)
	}
	return fmt.Sprintf("wrote %d bytes to %s", len(content), path), nil
}

