package main

import (
	"context"
	"log"
	"os"

	rootwebsitebuilder "seq-website-builder/agents/root-website-builder"
	"seq-website-builder/conf"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
)

func main() {
	conf := conf.LoadConfig()
	rootAgent, err := rootwebsitebuilder.SequentialAgent(conf.APIKey)
	if err != nil {
		log.Fatalf("Agent failed %v", err)
	}

	ctx := context.Background()
	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(rootAgent),
	}
	l := full.NewLauncher()
	if err = l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
