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

type PromptData struct {
	Prompt   string `json:"prompt"`
	Filename string `json:"filename"`
	Seed     int    `json:"seed"`
}

type PythonGenerationResult struct {
	Status      string `json:"status"`
	Output      string `json:"output"`
	PromptIndex int    `json:"prompt_index"`
}

// PythonOverallResult represents the overall JSON output from the Python script
type PythonOverallResult struct {
	OverallStatus string                   `json:"all_status"`
	Generations   []PythonGenerationResult `json:"generations"`
	Error         string                   `json:"error,omitempty"`
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

	var promptsData []PromptData
	for i, p := range promptConf.Prompts {
		prompt := fmt.Sprintf("%s %s", p, promptConf.StyleSuffix)
		filename := utils.SanitizeFilenameForImage(p, i+1)
		promptsData = append(promptsData, PromptData{
			Prompt:   prompt,
			Filename: filename,
			Seed:     cmdConf.Seed + i,
		})
	}

	promptsDataJSON, err := json.Marshal(promptsData)
	if err != nil {
		return fmt.Errorf("failed to marshal prompts data to JSON: %w", err)
	}

	args := []string{
		scriptPath,
		"--model", hfModel,
		"--negative-prompt", promptConf.NegativePrompt,
		"--width", strconv.Itoa(cmdConf.Resolution[0]),
		"--height", strconv.Itoa(cmdConf.Resolution[1]),
		"--steps", strconv.Itoa(cmdConf.Steps),
		"--guidance-scale", fmt.Sprintf("%.2f", cmdConf.GuidanceScale),
		"--output-dir", cmdConf.OutputDir,
		"--prompts-data", string(promptsDataJSON),
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

	var overallResult PythonOverallResult
	if err := json.Unmarshal(stdoutBuffer.Bytes(), &overallResult); err != nil {
		return fmt.Errorf("failed to parse python output: %v", err)
	}

	if overallResult.OverallStatus != "success" {
		return fmt.Errorf("overall generation failed: %s", overallResult.Error)
	}

	for _, res := range overallResult.Generations {
		if res.Status != "success" {
			// You might want to handle individual generation failures differently
			fmt.Printf("⚠ Generation for prompt index %d failed: %s\n", res.PromptIndex, res.Output)
		} else {
			fmt.Printf("    ✓ Saved to %s (Prompt %d) in %s\n", filepath.Base(res.Output), res.PromptIndex+1, time.Since(start))
		}
	}
	fmt.Printf(
		"\nTotal generation time for %d images: %s\n",
		len(promptConf.Prompts),
		time.Since(start),
	)

	return nil
}
