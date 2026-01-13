package codewriter

import (
	"context"
	"fmt"

	"seq-website-builder/conf"
	"seq-website-builder/tools"
	"seq-website-builder/utils"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/genai"
)

func CodeWriterAgent(c conf.Config, mdl string) (agent.Agent, error) {
	model, err := gemini.NewModel(context.Background(), mdl, &genai.ClientConfig{
		APIKey: c.APIKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create model: %v", err)
	}

	desc, err := utils.LoadInstructionsFile("./agents/code-writer/description.txt")
	if err != nil {
		return nil, err
	}
	instr, err := utils.LoadInstructionsFile("./agents/code-writer/instructions.txt")
	if err != nil {
		return nil, err
	}
	fileWriteTool, err := tools.NewFileWriteTool()
	if err != nil {
		return nil, fmt.Errorf("failed to create file write tool: %v", err)
	}

	agent, err := llmagent.New(llmagent.Config{
		Name:        "code_writer_agent",
		Model:       model,
		Description: desc,
		Instruction: instr,
		Tools: []tool.Tool{
			fileWriteTool,
		},
		// instruction and tools will be added next
	})
	if err != nil {
		return nil, fmt.Errorf("error constructing the designer agent %v", err)
	}

	return agent, nil
}
