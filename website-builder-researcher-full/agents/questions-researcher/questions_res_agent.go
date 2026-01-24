package questionsreasearcher

import (
	"context"
	"fmt"

	"seq-website-builder/utils"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/agent/workflowagents/parallelagent"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/geminitool"
	"google.golang.org/genai"
)

func QuestionResearcher(APIKey string) (agent.Agent, error) {
	model, err := gemini.NewModel(context.Background(), "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: APIKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create model: %v", err)
	}
	fmt.Println("Question researcher agent started")

	desc, err := utils.LoadInstructionsFile("./agents/questions-researcher/description.txt")
	if err != nil {
		return nil, err
	}
	instr, err := utils.LoadInstructionsFile("./agents/questions-researcher/instructions.txt")
	if err != nil {
		return nil, err
	}

	researcher1, err := llmagent.New(llmagent.Config{
		Name:  "QuestionResearcher1",
		Model: model,
		Instruction: fmt.Sprintf(
			"You are assigned to answer QUESTION NUMBER 1 only.\n\n%s",
			instr,
		),
		Description: fmt.Sprintf(
			"%s This agent specifically handles question #1.",
			desc,
		),
		OutputKey: "question_1_research_output",
		Tools: []tool.Tool{
			geminitool.GoogleSearch{},
		},
	})
	if err != nil {
		return nil, err
	}
	researcher2, err := llmagent.New(llmagent.Config{
		Name:  "QuestionResearcher2",
		Model: model,
		Instruction: fmt.Sprintf(
			"You are assigned to answer QUESTION NUMBER 2 only.\n\n%s",
			instr,
		),
		Description: fmt.Sprintf(
			"%s This agent specifically handles question #2.",
			desc,
		),
		OutputKey: "question_2_research_output",
		Tools: []tool.Tool{
			geminitool.GoogleSearch{},
		},
	})
	if err != nil {
		return nil, err
	}

	parallelResearchAgent, err := parallelagent.New(parallelagent.Config{
		AgentConfig: agent.Config{
			Name:        "ParallelQuestionsResearchAgent",
			Description: "Runs five question research agents in parallel to research and answer all five questions simultaneously.",
			SubAgents: []agent.Agent{
				researcher1,
				researcher2,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create parallel agent: %v", err)
	}
	fmt.Println("Question Researched agent finished")
	return parallelResearchAgent, nil
}
