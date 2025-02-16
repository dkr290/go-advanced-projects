package helpers

import (
	"context"
	"os/exec"
)

// Helper function to execute prompt, it could be that many different functions will execute different commands
func LlamaCommandPrompt(
	ctx context.Context,
	modelPath, prompt, llamaPath string,
) (*exec.Cmd, error) {
	cmd := exec.CommandContext(ctx, llamaPath, "-m", modelPath, "-p", prompt)

	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return cmd, nil
}
