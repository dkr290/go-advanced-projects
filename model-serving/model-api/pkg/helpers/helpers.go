package helpers

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"time"
)

func ConvertMultiSafetensors(tempDir, outputDir string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	// Get sorted list of parts
	files, err := os.ReadDir(tempDir)
	if err != nil {
		return err
	}

	// Sort files numerically
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	// Build input pattern (assumes sequential numbering)
	pattern := filepath.Join(tempDir, "part-*.safetensors")

	cmd := exec.CommandContext(ctx,
		"llama.cpp/build/bin/llama-gguf",
		"--input", pattern, // Use wildcard pattern
		"--output", filepath.Join(outputDir, "model.gguf"),
		"--ctx", "4096",
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w (stderr: %s)", err, stderr.String())
	}

	return nil
}
