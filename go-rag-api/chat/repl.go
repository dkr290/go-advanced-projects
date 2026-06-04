package chat

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/dkr290/go-advanced-projects/go-rag-api/llm"
)

type Options struct {
	SystemPromptFile string
}

func RunREPL(ctx context.Context, client *llm.Client, opts Options) error {
	in := bufio.NewScanner(os.Stdin)
	in.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	history, err := seedHistory(opts.SystemPromptFile)
	if err != nil {
		return err
	}
	fmt.Println("Chat session started. Type Q/q to quit")
	for {
		fmt.Print("\n> ")
		if !in.Scan(){
			if err := in.Err();err != nil {
				return err
			}
		}
	}
}

func seedHistory(path string) ([]llm.Message, error) {
	if path == "" {
		return nil, nil
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read system prompt: %w", err)
	}

	content := strings.TrimSpace(string(data))
	if content == "" {
		return nil, nil
	}

 msg := []llm.Message {

		{Role: "system", Content: content},
}
return  msg,nil
}
