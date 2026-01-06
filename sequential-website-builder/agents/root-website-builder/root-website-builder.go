package rootwebsitebuilder

import (
	"fmt"

	codewriter "seq-website-builder/agents/code-writer"
	"seq-website-builder/agents/designer"
	requirementswriter "seq-website-builder/agents/requirements-writer"
	"seq-website-builder/utils"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/workflowagents/sequentialagent"
)

func SequentialAgent(apiKey string) (agent.Agent, error) {
	codeWriterAgent, err := codewriter.CodeWriterAgent(apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create code writer agent: %v", err)
	}

	designerAgent, err := designer.Designer(apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create code writer agent: %v", err)
	}
	requrenmentsWriterAgent, err := requirementswriter.Writer(apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create code writer agent: %v", err)
	}

	desc, err := utils.LoadInstructionsFile("./description.txt")
	if err != nil {
		return nil, err
	}

	codePipelineAgent, err := sequentialagent.New(sequentialagent.Config{
		AgentConfig: agent.Config{
			Name:        "root_website_builder_agent",
			Description: desc,
			SubAgents: []agent.Agent{
				requrenmentsWriterAgent,
				designerAgent,
				codeWriterAgent,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create sequential agent: %v", err)
	}
	return codePipelineAgent, nil
}
