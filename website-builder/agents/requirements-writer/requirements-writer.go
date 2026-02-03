// Package requirementswriter for requirements_writer agent
package requirementswriter

import (
	"context"
	"fmt"

	"website-builder/logs"
	"website-builder/utils"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/genai"
)

func Writer(APIKey, m string, lloger *logs.Logger) (agent.Agent, error) {
	model, err := gemini.NewModel(context.Background(), m, &genai.ClientConfig{
		APIKey: APIKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create model: %v", err)
	}

	desc, err := utils.LoadInstructionsFile("./agents/requirements-writer/description.txt")
	if err != nil {
		return nil, err
	}
	instr, err := utils.LoadInstructionsFile("./agents/requirements-writer/instructions.txt")
	if err != nil {
		return nil, err
	}

	agent, err := llmagent.New(llmagent.Config{
		Name:        "requirements_writer_agent",
		Model:       model,
		Description: desc,
		Instruction: instr,
		OutputKey:   "requirements_writer_output",
		// instruction and tools will be added next
	})
	if err != nil {
		return nil, fmt.Errorf("error constructing the agent %v", err)
	}

	return agent, nil
}
