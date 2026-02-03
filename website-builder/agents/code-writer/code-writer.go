package codewriter

import (
	"context"
	"fmt"

	"website-builder/logs"
	"website-builder/tools"
	"website-builder/utils"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/genai"
)

func CodeWriterAgent(APIKey, m string, lloger *logs.Logger) (agent.Agent, error) {
	model, err := gemini.NewModel(context.Background(), m, &genai.ClientConfig{
		APIKey: APIKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create model: %v", err)
	}
	lloger.Logging.Debugf("Loading the model %s", model.Name())

	lloger.Logging.Debugln("Loading the description file")
	desc, err := utils.LoadInstructionsFile("./agents/code-writer/description.txt")
	if err != nil {
		return nil, err
	}

	lloger.Logging.Debugln("Loading the instruction file")

	instr, err := utils.LoadInstructionsFile("./agents/code-writer/instructions.txt")
	if err != nil {
		return nil, err
	}
	fileWriteTool, err := tools.NewFileWriteTool()
	if err != nil {
		return nil, fmt.Errorf("failed to create file write tool: %v", err)
	}

	lloger.Logging.Debugln("Running the code writer agent")

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
