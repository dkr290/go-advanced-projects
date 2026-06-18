package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/example/go-code-agent/pkg/ai"
)

func renderSpinner(stopCh <-chan struct{}) {
	frames := []string{"|", "/", "-", "\\"}
	i := 0

	for {
		select {
		case <-stopCh:
			// wipe the line clear when thinking finishes
			fmt.Print("\r\033[K")
			return
		default:
			fmt.Printf("\r\033[35m%s\033[0m AI is reviewing files and thinking...", frames[i%4])
			i++
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func main() {
	ctx := context.Background()

	modelName := "deepseek-coder"
	localAIURL := "http://localhost:8080"
	fmt.Println("\033[36m=== LocalAI CLI Agent Ready ===\033[0m")
	fmt.Printf("Model: %s at %s\n", modelName, localAIURL)
	fmt.Println("Type your instructions below (type 'exit' to quit):")
	agent, err := ai.SetupAgent(modelName, localAIURL)
	if err != nil {
		log.Fatalf("Agent build sequence crashed: %v", err)
	}
	defer agent.Cleanup(ctx)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\n\033[32m[CodeAgent Ask]>\033[0m ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())

		if input == "" {
			continue
		}
		if strings.ToLower(input) == "exit" {
			break
		}

		stopSpinner := make(chan struct{})
		go renderSpinner(stopSpinner)

		// per-request timeout
		reqCtx, reqCancel := context.WithTimeout(ctx, 2*time.Minute)

		// Use streaming for real-time response
		response, err := agent.Run(reqCtx, input)
		reqCancel()

		close(stopSpinner)
		time.Sleep(50 * time.Millisecond)

		if err != nil {
			fmt.Printf("\n\033[31m[Error]:\033[0m %v\n", err)
			continue
		}

		fmt.Println("\n\033[34m[Response]:\033[0m")
		fmt.Println(response.Content)

	}
	if err := scanner.Err(); err != nil {
		log.Printf("scanner error: %v", err)
	}
}
