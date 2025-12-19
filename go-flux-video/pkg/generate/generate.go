// Package generate - is where the actual image generation happened
package generate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gfluxgo/pkg/config"
	"gfluxgo/pkg/utils"
)

type PythonResult struct {
	Status string `json:"status"`
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
}

// GenerateWithPython calls Python for FLUX GGUF generation
func GenerateWithPython(
	cmdConf config.Config,
	promptConf config.PromptConfig,
	modelPath, loraDir string,
) error {
	scriptPath := filepath.Join("scripts", "python_generate.py")

	// Check if script exists
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return fmt.Errorf("python script not found at %s", scriptPath)
	}

	if err := os.MkdirAll(cmdConf.OutputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output dir: %w", err)
	}

	// Determine if modelPath is GGUF or HuggingFace ID
	var hfModel, ggufPath string
	if strings.HasSuffix(strings.ToLower(modelPath), ".gguf") {
		ggufPath = modelPath
		hfModel = "black-forest-labs/FLUX.1-dev"
	} else {
		hfModel = modelPath
	}

	// Find LoRA file if loraDir is provided
	var loraFilePath string
	if loraDir != "" {
		entries, err := os.ReadDir(loraDir)
		if err == nil {
			for _, e := range entries {
				if strings.HasSuffix(e.Name(), ".safetensors") {
					loraFilePath = filepath.Join(loraDir, e.Name())
					break
				}
			}
		}
	}

	fmt.Printf("\nStarting FLUX generation for %d images...\n", len(promptConf.Prompts))

	for i, p := range promptConf.Prompts {
		prompt := fmt.Sprintf("%s, %s", p, promptConf.StyleSuffix)
		filename := utils.SanitizeFilenameForImage(p, i+1)
		outputPath := filepath.Join(cmdConf.OutputDir, filename)

		fmt.Printf("[%d/%d] Generating: %s\n", i+1, len(promptConf.Prompts), p)

		args := []string{
			scriptPath,
			"--model", hfModel,
			"--prompt", prompt,
			"--negative-prompt", promptConf.NegativePrompt,
			"--width", strconv.Itoa(cmdConf.Resolution[0]),
			"--height", strconv.Itoa(cmdConf.Resolution[1]),
			"--steps", strconv.Itoa(cmdConf.Steps),
			"--guidance-scale", fmt.Sprintf("%.2f", cmdConf.GuidanceScale),
			"--seed", strconv.Itoa(cmdConf.Seed + i),
			"--output", outputPath,
		}

		if ggufPath != "" {
			args = append(args, "--gguf", ggufPath)
		}

		// Pass both LoRA repo ID and file path
		if cmdConf.LoraRepo != "" && loraFilePath != "" {
			args = append(args, "--lora", cmdConf.LoraRepo)
			args = append(args, "--lora-file", loraFilePath)
		}
		// Add low VRAM flag if enabled
		if cmdConf.LowVRAM {
			args = append(args, "--low-vram")
		}

		start := time.Now()

		cmd := exec.Command("python3", args...)

		var stdout bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = io.MultiWriter(os.Stderr, &stderr) // Print AND capture

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("python failed: %w\nstderr: %s", err, stderr.String())
		}

		var result PythonResult
		if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
			return fmt.Errorf(
				"parse error: %w\nstdout: %s\nstderr: %s",
				err,
				stdout.String(),
				stderr.String(),
			)
		}

		if result.Status != "success" {
			return fmt.Errorf("generation failed: %s\nstderr: %s", result.Error, stderr.String())
		}

		fmt.Printf("    ✓ Saved to %s in %s\n", filename, time.Since(start))
	}

	return nil
}
