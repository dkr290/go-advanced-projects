// Package questiongenerator
package questiongenerator

import (
	"context"
	"fmt"

	"seq-website-builder/conf"
	"seq-website-builder/utils"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/geminitool"
	"google.golang.org/genai"
)

func QuestionGenerator(c conf.Config, mdl string) (agent.Agent, error) {
	model, err := gemini.NewModel(context.Background(), mdl, &genai.ClientConfig{
		APIKey: c.APIKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create model: %v", err)
	}

	desc, err := utils.LoadInstructionsFile("./agents/designer/description.txt")
	if err != nil {
		return nil, err
	}
	instr, err := utils.LoadInstructionsFile("./agents/designer/instructions.txt")
	if err != nil {
		return nil, err
	}

	agent, err := llmagent.New(llmagent.Config{
		Name:        "questions_generator_agent",
		Model:       model,
		Description: desc,
		Instruction: instr,
		OutputKey:   "questions_generator_output",
		Tools: []tool.Tool{
			geminitool.GoogleSearch{},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error constructing the designer agent %v", err)
	}

	return agent, nil
}
