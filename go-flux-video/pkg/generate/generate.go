// Package generate - is where the actual image generation happened
package generate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gfluxgo/pkg/config"
	"gfluxgo/pkg/utils"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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
		hfModel = cmdConf.HfModelID
	} else {
		hfModel = cmdConf.HfModelID
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

		fmt.Printf("[%d/%d] Generating: %s\n", i+1, len(promptConf.Prompts), prompt)

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

		stdoutPipe, err := cmd.StdoutPipe()
		if err != nil {
			return fmt.Errorf("failed to create stdout pipe: %w", err)
		}

		// Stream stderr directly to terminal for progress
		cmd.Stderr = os.Stderr

		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to start python: %w", err)
		}

		var stdoutBuffer bytes.Buffer
		// MultiWriter writes to both buffer and terminal
		stdoutMulti := io.MultiWriter(&stdoutBuffer, os.Stdout)
		if _, err := io.Copy(stdoutMulti, stdoutPipe); err != nil {
			return fmt.Errorf("failed to read python output: %w", err)
		}

		if err := cmd.Wait(); err != nil {
			return fmt.Errorf("python execution failed: %v", err)
		}

		var result PythonResult
		if err := json.Unmarshal(stdoutBuffer.Bytes(), &result); err != nil {
			return fmt.Errorf("failed to parse python output: %v", err)
		}
		if result.Status != "success" {
			return fmt.Errorf("generation failed: %s", result.Error)
		}

		fmt.Printf("    ✓ Saved to %s in %s\n", filename, time.Since(start))
	}

	return nil
}
