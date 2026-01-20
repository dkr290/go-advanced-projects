package rootwebsitebuilder

import (
	"fmt"

	requirementswriter "seq-website-builder/agents/requirements-writer"
	"seq-website-builder/conf"
	"seq-website-builder/utils"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/workflowagents/sequentialagent"
)

type WebSiteBuilderConfig struct {
	config conf.Config
	l      utils.Logger
}

func NewBuilderAgent(
	config conf.Config, log utils.Logger,
) *WebSiteBuilderConfig {
	return &WebSiteBuilderConfig{
		config: config,
		l:      log,
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
	// Validate we have models
	if len(w.config.Models) == 0 {
		return nil, fmt.Errorf("no models configured")
	}

	w.l.Logging.Debugf("All models selected %s", w.config.Models)

	// questiongeneratorAgent, err := questiongenerator.QuestionGenerator(
	// 	w.config,
	// 	getModel(0),
	// 	w.l,
	// )
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create question generator agent: %v", err)
	// }
	w.l.Logging.Debugf("The models for question generation agent: %v", getModel(0))
	w.l.Logging.Debugf("The api key : %v", w.config.APIKey)

	// questionResearcherAgent, err := questionsreasearcher.QuestionResearcher(
	// 	w.config,
	// 	getModel(0),
	// )
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create question researcher agent: %v", err)
	// }
	//
	// w.l.Logging.Debugf("The model for question research agent: %v", getModel(0))
	//
	// queryGeneratorAgent, err := querygenerator.QueryGenerator(w.config, getModel(0))
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create query generator agent: %v", err)
	// }

	w.l.Logging.Debugf("The model for query generator agent: %v", getModel(0))

	// codeWriterAgent, err := codewriter.CodeWriterAgent(w.config, getModel(1))
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create code writer agent: %v", err)
	// }

	w.l.Logging.Debugf("The model for Code Writer agent %s", getModel(1))

	// designerAgent, err := designer.Designer(w.config, getModel(2))
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create designer writer agent: %v", err)
	// }

	w.l.Logging.Debugf("The model for Designer agent %s", getModel(2))

	requirenmentsWriterAgent, err := requirementswriter.Writer(w.config, getModel(3))
	if err != nil {
		return nil, fmt.Errorf("failed to create requirements writer agent: %v", err)
	}
	w.l.Logging.Debugf("The model for Requirements Writer agent %s", getModel(3))

	desc, err := utils.LoadInstructionsFile("./agents/root-website-builder/description.txt")
	if err != nil {
		return nil, err
	}

	codePipelineAgent, err := sequentialagent.New(sequentialagent.Config{
		AgentConfig: agent.Config{
			Name:        "root_website_builder_agent",
			Description: desc,
			SubAgents: []agent.Agent{
				// questiongeneratorAgent,
				// questionResearcherAgent,
				// queryGeneratorAgent,
				requirenmentsWriterAgent,
				// designerAgent,
				// codeWriterAgent,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create sequential agent: %v", err)
	}
	return codePipelineAgent, nil
}
