package app

import (
	"context"

	"github.com/dkr290/go-advanced-projects/go-rag-api/chat"
	"github.com/dkr290/go-advanced-projects/go-rag-api/config"
	"github.com/dkr290/go-advanced-projects/go-rag-api/llm"
)


func Run(ctx context.Context,cfg config.Config) error {

	client := llm.New(cfg)
	return chat.RunREPL(ctx,client,chat.Options{

		SystemPromptFile: cfg.SystemPromptFile,


	})

}
