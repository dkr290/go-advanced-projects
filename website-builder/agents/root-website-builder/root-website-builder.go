package rootwebsitebuilder

import (
	"context"
	"fmt"

	"website-builder/agents/designer"
	"website-builder/conf"
	"website-builder/logs"
	"website-builder/utils"

	codewriter "website-builder/agents/code-writer"

	requirementswriter "website-builder/agents/requirements-writer"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/workflowagents/sequentialagent"
	"google.golang.org/genai"
)

type RootWebsite struct {
	conf  conf.Config
	model string
}

func NewRootBulder(cn conf.Config) (*RootWebsite, error) {
	for _, model := range cn.Models {
		err := tryModel(cn.APIKey, model)
		if err == nil {
			fmt.Println("Using model", model)
			return &RootWebsite{
				conf:  cn,
				model: model,
			}, nil
		}
		// Log and try next model
		fmt.Printf("Model %s failed: %v, trying next...\n", model, err)
	}
	return nil, fmt.Errorf("all models failed, tried: %v", cn.Models)
}

// tryModel attempts to create an agent with the given model to verify it works
func tryModel(apiKey, model string) error {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	resp, err := client.Models.GenerateContent(ctx, model, genai.Text("Hello!"), nil)
	if err != nil {
		return fmt.Errorf("generation failed: %w", err)
	}

	// 4. Handle the Response
	if len(resp.Candidates) > 0 {
		fmt.Println("Model Response:", resp.Candidates[0].Content.Parts[0])
	}

	return nil
}

func (r *RootWebsite) SequentialAgent() (agent.Agent, error) {
	llogger := logs.Init(r.conf.Debug)

	codeWriterAgent, err := codewriter.CodeWriterAgent(r.conf.APIKey, r.model, llogger)
	if err != nil {
		return nil, fmt.Errorf("failed to create code writer agent: %v", err)
	}

	designerAgent, err := designer.Designer(r.conf.APIKey, r.model, llogger)
	if err != nil {
		return nil, fmt.Errorf("failed to create designer writer agent: %v", err)
	}
	requirenmentsWriterAgent, err := requirementswriter.Writer(r.conf.APIKey, r.model, llogger)
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
