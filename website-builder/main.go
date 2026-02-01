package main

import (
	"context"
	"log"
	"os"

	rootwebsitebuilder "website-builder/agents/root-website-builder"
	"website-builder/conf"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
)

func main() {
	conf := conf.LoadConfig()
	rootWebsiteBuilder, err := rootwebsitebuilder.NewRootBulder(*conf)
	if err != nil {
		log.Fatalf("Unable to call root website buidler %v", err)
	}
	rootAgent, err := rootWebsiteBuilder.SequentialAgent()
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
