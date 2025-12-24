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

// GenerateWithPythonQwen calls Python for Qwen-Image-Edit text-to-image generation
func GenerateWithPythonQwen(
	cmdConf config.Config,
	promptConf config.PromptConfig,
	modelPath, loraDir string,
) error {
	scriptPath := filepath.Join("scripts", "qwen_img2img.py")

	// Check if script exists
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return fmt.Errorf("qwen python script not found at %s", scriptPath)
	}

	if err := os.MkdirAll(cmdConf.OutputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output dir: %w", err)
	}

	// Use HuggingFace model ID directly for Qwen
	hfModel := cmdConf.HfModelID

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

	fmt.Printf("\nStarting Qwen-Image-Edit generation for %d images...\n", len(promptConf.Prompts))

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
		"--strength", fmt.Sprintf("%.2f", cmdConf.Strength),
		"--output-dir", cmdConf.OutputDir,
		"--prompts-data", string(promptsDataJSON),
		"--num-images", strconv.Itoa(cmdConf.QwenNumImages),
	}

	// Add LoRA if found
	if loraFilePath != "" {
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

	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start python: %w", err)
	}

	var stdoutBuffer bytes.Buffer
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
			fmt.Printf("⚠ Generation for prompt index %d failed: %s\n", res.PromptIndex, res.Output)
		} else {
			fmt.Printf("    ✓ Saved to %s (Prompt %d) in %s\n", filepath.Base(res.Output), res.PromptIndex+1, time.Since(start))
		}
	}

	fmt.Printf(
		"\nTotal Qwen generation time for %d images: %s\n",
		len(promptConf.Prompts),
		time.Since(start),
	)

	return nil
}

// GenerateImg2ImgWithPythonQwen calls Python for Qwen-Image-Edit image-to-image generation
func GenerateImg2ImgWithPythonQwen(
	cmdConf config.Config,
	promptConf config.PromptConfig,
	modelPath, loraDir string,
	inputImagesDir string, // Directory containing input images
) error {
	scriptPath := filepath.Join("scripts", "qwen_img2img.py")

	// Check if script exists
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return fmt.Errorf("qwen python img2img script not found at %s", scriptPath)
	}

	if err := os.MkdirAll(cmdConf.OutputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output dir: %w", err)
	}

	// Use HuggingFace model ID directly for Qwen
	hfModel := cmdConf.HfModelID

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

	fmt.Printf(
		"\nStarting Qwen-Image-Edit img2img generation for %d images...\n",
		len(promptConf.Prompts),
	)

	// Read input images from directory
	inputImages, err := os.ReadDir(inputImagesDir)
	if err != nil {
		return fmt.Errorf("failed to read input images directory: %w", err)
	}

	var promptsData []PromptData
	imageIdx := 0
	for i, p := range promptConf.Prompts {
		// Match prompt with input image
		var inputImagePath string
		if imageIdx < len(inputImages) && !inputImages[imageIdx].IsDir() {
			inputImagePath = filepath.Join(inputImagesDir, inputImages[imageIdx].Name())
			imageIdx++
		} else {
			return fmt.Errorf("not enough input images for prompt %d", i)
		}

		prompt := fmt.Sprintf("%s %s", p, promptConf.StyleSuffix)
		filename := utils.SanitizeFilenameForImage(p, i+1)
		promptsData = append(promptsData, PromptData{
			Prompt:     prompt,
			Filename:   filename,
			Seed:       cmdConf.Seed + i,
			InputImage: inputImagePath,
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
		"--strength", fmt.Sprintf("%.2f", cmdConf.Strength),
		"--output-dir", cmdConf.OutputDir,
		"--prompts-data", string(promptsDataJSON),
		"--num-images", strconv.Itoa(cmdConf.QwenNumImages),
	}

	// Add LoRA if found
	if loraFilePath != "" {
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

	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start python: %w", err)
	}

	var stdoutBuffer bytes.Buffer
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
			fmt.Printf("⚠ Generation for prompt index %d failed: %s\n", res.PromptIndex, res.Output)
		} else {
			fmt.Printf("    ✓ Saved to %s (Prompt %d) in %s\n", filepath.Base(res.Output), res.PromptIndex+1, time.Since(start))
		}
	}

	fmt.Printf(
		"\nTotal Qwen img2img generation time for %d images: %s\n",
		len(promptConf.Prompts),
		time.Since(start),
	)

	return nil
}
