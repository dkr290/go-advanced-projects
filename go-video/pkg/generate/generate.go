// Package generate - is where the actual image generation happened
package generate

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"time"

	"gfluxgo/pkg/config"
	"gfluxgo/pkg/utils"

	"github.com/binozo/gostablediffusion/pkg/sd"
)

func Generate(
	cmdConf config.Config,
	promptConf config.PromptConfig,
	modelPath, loraDir string,
) error {
	sdBuilder := sd.New().SetModel(modelPath)
	// Only set the LoRA dir if it contains at least one .safetensors file
	if cmdConf.LoraURL != "" {
		fmt.Println("Setting the lora directory")
		sdBuilder.SetLoRaDir(loraDir)
	}

	ctx, err := sdBuilder.Load()
	if err != nil {
		return fmt.Errorf("CGO model load error (GGUF/driver mismatch?): %w", err)
	}
	// Validate context is not nil
	if ctx == nil {
		return fmt.Errorf("context is nil after loading model")
	}

	defer ctx.Free() // Releases C++ resources

	// --- 4. Generation Loop ---
	outputDir := cmdConf.OutputDir
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("failed creating output dir %q: %w", outputDir, err)
	}

	fmt.Printf("\nStarting character sheet generation for %d poses...\n", len(promptConf.Prompts))

	// Validate resolution
	if len(cmdConf.Resolution) != 2 {
		return fmt.Errorf(
			"invalid resolution format, expected [width, height], got %v",
			cmdConf.Resolution,
		)
	}
	width, height := cmdConf.Resolution[0], cmdConf.Resolution[1]
	if width <= 0 || height <= 0 {
		return fmt.Errorf("invalid resolution: width=%d, height=%d", width, height)
	}

	for i, p := range promptConf.Prompts {
		// Construct the final prompt
		prompt := fmt.Sprintf("%s, %s", p, promptConf.StyleSuffix)

		fmt.Printf("[%d/%d] Generating: %s\n", i+1, len(promptConf.Prompts), p)

		// Create the generation parameters
		params := sd.NewImageGenerationParams()
		params.Width = width
		params.Height = height
		params.Prompt = prompt
		params.NegativePrompt = promptConf.NegativePrompt
		params.SampleSteps = cmdConf.Steps
		params.Guidance.TxtCfg = cmdConf.GuidanceScale
		params.Seed = int64(cmdConf.Seed + i)

		// Validate parameters
		if params.Prompt == "" {
			fmt.Printf("Warning: Empty prompt for pose %d, skipping\n", i+1)
			continue
		}

		start := time.Now()

		// CGO CALL: The generation happens here
		fmt.Printf(
			"    Calling CGO GenerateImage with params: width=%d, height=%d, steps=%d, seed=%d\n",
			params.Width,
			params.Height,
			params.SampleSteps,
			params.Seed,
		)

		// CGO CALL: The generation happens here
		imageData := ctx.GenerateImage(params)
		img := imageData.Image()
		if img == nil {
			imageData.Free()
			return fmt.Errorf(
				"GenerateImage returned a nil image (model context may not be initialized)",
			)
		}

		// Save Image
		safeLabel := utils.SanitizeFilename(p)
		filename := fmt.Sprintf("%02d_%s_s%d.png", i+1, safeLabel, params.Seed)
		outputPath := filepath.Join(outputDir, filename)

		targetFile, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE, 0o600)
		if err != nil {
			imageData.Free()
			return fmt.Errorf("creating file %s: %w", outputPath, err)
		}

		if err = png.Encode(targetFile, imageData.Image()); err != nil {
			targetFile.Close()
			imageData.Free()
			return fmt.Errorf("encoding PNG %s: %w", outputPath, err)
		}
		targetFile.Close()
		imageData.Free()

		fmt.Printf("    ✓ Saved to %s in %s\n", filename, time.Since(start))
	}
	return nil
}
