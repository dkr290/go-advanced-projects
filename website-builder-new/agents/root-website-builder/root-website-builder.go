package rootwebsitebuilder

import (
	"fmt"

	codewriter "seq-website-builder/agents/code-writer"
	"seq-website-builder/agents/designer"
	querygenerator "seq-website-builder/agents/query-generator"
	questiongenerator "seq-website-builder/agents/question-generator"
	questionsreasearcher "seq-website-builder/agents/questions-researcher"
	requirementswriter "seq-website-builder/agents/requirements-writer"
	"seq-website-builder/utils"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/workflowagents/sequentialagent"
)

func SequentialAgent(apiKey string) (agent.Agent, error) {
	questiongeneratorAgent, err := questiongenerator.QuestionGenerator(apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create question generator agent: %v", err)
	}

	questionResearcherAgent, err := questionsreasearcher.QuestionResearcher(apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create question researcher agent: %v", err)
	}

	queryGeneratorAgent, err := querygenerator.QueryGenerator(apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create query generator agent: %v", err)
	}
	requirenmentsWriterAgent, err := requirementswriter.Writer(apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create requirements writer agent: %v", err)
	}

	codeWriterAgent, err := codewriter.CodeWriterAgent(apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create code writer agent: %v", err)
	}

	designerAgent, err := designer.Designer(apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create designer writer agent: %v", err)
	}

	desc, err := utils.LoadInstructionsFile("./agents/root-website-builder/description.txt")
	if err != nil {
		return nil, err
	}

	codePipelineAgent, err := sequentialagent.New(sequentialagent.Config{
		AgentConfig: agent.Config{
			Name:        "root_website_builder_agent",
			Description: desc,
			SubAgents: []agent.Agent{
				questiongeneratorAgent,
				questionResearcherAgent,
				queryGeneratorAgent,
				requirenmentsWriterAgent,
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
