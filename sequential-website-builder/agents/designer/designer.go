// Package designer
package designer

import (
	"context"
	"fmt"

	"seq-website-builder/utils"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/genai"
)

func Designer(APIKey string) (agent.Agent, error) {
	model, err := gemini.NewModel(context.Background(), "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: APIKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create model: %v", err)
	}

	desc, err := utils.LoadInstructionsFile("./description.txt")
	if err != nil {
		return nil, err
	}
	instr, err := utils.LoadInstructionsFile("./instructions.txt")
	if err != nil {
		return nil, err
	}

	agent, err := llmagent.New(llmagent.Config{
		Name:        "designer_agent",
		Model:       model,
		Description: desc,
		Instruction: instr,
		OutputKey:   "designer_output",
		// instruction and tools will be added next
	})
	if err != nil {
		return nil, fmt.Errorf("error constructing the designer agent %v", err)
	}

	return agent, nil
}
