// Package designer
package designer

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

func Designer(APIKey, m string, lloger *logs.Logger) (agent.Agent, error) {
	model, err := gemini.NewModel(context.Background(), m, &genai.ClientConfig{
		APIKey: APIKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create model: %v", err)
	}
	lloger.Logging.Debugf("Loaded the model %s", model.Name())

	lloger.Logging.Debugln("Loading the description file")

	desc, err := utils.LoadInstructionsFile("./agents/designer/description.txt")
	if err != nil {
		return nil, err
	}

	lloger.Logging.Debugln("Loading the instruction file")

	instr, err := utils.LoadInstructionsFile("./agents/designer/instructions.txt")
	if err != nil {
		return nil, err
	}
	lloger.Logging.Debugln("Running the designer agent")

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
