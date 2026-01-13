package rootwebsitebuilder

import (
	"fmt"

	codewriter "seq-website-builder/agents/code-writer"
	"seq-website-builder/agents/designer"
	querygenerator "seq-website-builder/agents/query-generator"
	questiongenerator "seq-website-builder/agents/question-generator"
	questionsreasearcher "seq-website-builder/agents/questions-reasearcher"
	requirementswriter "seq-website-builder/agents/requirements-writer"
	"seq-website-builder/conf"
	"seq-website-builder/utils"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/workflowagents/sequentialagent"
)

type WebSiteBuilderConfig struct {
	config conf.Config
}

func NewBuilderAgent(
	config conf.Config,
) *WebSiteBuilderConfig {
	return &WebSiteBuilderConfig{
		config: config,
	}
}

func (w *WebSiteBuilderConfig) SequentialAgent() (agent.Agent, error) {
	// Helper to safely get model index, fallback to n-1
	getModel := func(idx int) string {
		for i := idx; i >= 0; i-- {
			if len(w.config.Models) > i {
				return w.config.Models[i]
			}
		}
		return ""
	}

	questiongeneratorAgent, err := questiongenerator.QuestionGenerator(
		w.config,
		getModel(0),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create question generator agent: %v", err)
	}
	questionResearcherAgent, err := questionsreasearcher.QuestionResearcher(
		w.config,
		getModel(0),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create question researcher agent: %v", err)
	}

	queryGeneratorAgent, err := querygenerator.QueryGenerator(w.config, getModel(0))
	if err != nil {
		return nil, fmt.Errorf("failed to create query generator agent: %v", err)
	}

	codeWriterAgent, err := codewriter.CodeWriterAgent(w.config, getModel(1))
	if err != nil {
		return nil, fmt.Errorf("failed to create code writer agent: %v", err)
	}

	designerAgent, err := designer.Designer(w.config, getModel(2))
	if err != nil {
		return nil, fmt.Errorf("failed to create designer writer agent: %v", err)
	}
	requirenmentsWriterAgent, err := requirementswriter.Writer(w.config, getModel(3))
	if err != nil {
		return nil, fmt.Errorf("failed to create requirements writer agent: %v", err)
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
