package main

import (
	"context"
	"log"
	"os"
	"strconv"

	rootwebsitebuilder "seq-website-builder/agents/root-website-builder"
	"seq-website-builder/conf"
	"seq-website-builder/utils"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
)

func main() {
	conf := conf.LoadConfig()

	b, err := strconv.ParseBool(conf.Debug)
	if err != nil {
		log.Fatal(err)
	}

	llogger := utils.Init(b)
	wb := rootwebsitebuilder.NewBuilderAgent(*conf, *llogger)
	rootAgent, err := wb.SequentialAgent()
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
