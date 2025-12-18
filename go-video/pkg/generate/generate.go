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
	sdBuilder.SetLoRaDir(loraDir)
	ctx, err := sdBuilder.Load()
	if err != nil {
		fmt.Printf("FATAL CGO Model Load Error (Is your GGUF file correct?): %v\n", err)
		os.Exit(1)
	}
	defer ctx.Free() // Releases C++ resources

	// --- 4. Generation Loop ---
	outputDir := cmdConf.OutputDir
	os.MkdirAll(outputDir, 0o755)

	fmt.Printf("\nStarting character sheet generation for %d poses...\n", len(promptConf.Prompts))

	for i, p := range promptConf.Prompts {
		// Construct the final prompt
		prompt := fmt.Sprintf("%s, %s", p, promptConf.StyleSuffix)

		fmt.Printf("[%d/%d] Generating: %s\n", i+1, len(promptConf.Prompts), p)

		// Create the generation parameters
		params := sd.NewImageGenerationParams()
		params.Width = cmdConf.Resolution[0]
		params.Height = cmdConf.Resolution[1]
		params.Prompt = prompt
		params.NegativePrompt = promptConf.NegativePrompt
		params.SampleSteps = cmdConf.Steps
		params.Guidance.TxtCfg = cmdConf.GuidanceScale
		params.Seed = int64(cmdConf.Seed + i)

		start := time.Now()
		// CGO CALL: The generation happens here
		imageData := ctx.GenerateImage(params)
		// Save Image
		safeLabel := utils.SanitizeFilename(p)
		filename := fmt.Sprintf("%02d_%s_s%d.png", i+1, safeLabel, params.Seed)
		outputPath := filepath.Join(outputDir, filename)

		targetFile, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE, 0o600)
		if err != nil {
			imageData.Free()
			fmt.Printf("Error creating file %s: %v\n", outputPath, err)
			os.Exit(1)
		}

		if err = png.Encode(targetFile, imageData.Image()); err != nil {
			targetFile.Close()
			imageData.Free()
			fmt.Printf("Error encoding PNG %s: %v\n", outputPath, err)
			os.Exit(1)
		}
		targetFile.Close()
		imageData.Free()

		fmt.Printf("    ✓ Saved to %s in %s\n", filename, time.Since(start))
	}
	return nil
}
