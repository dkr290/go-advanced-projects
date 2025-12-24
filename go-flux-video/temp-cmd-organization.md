> Context:
>
> -  <buf>pkg/config/config.go</buf>

if you read the buffer #{buffer:config.go} How to organize the command line like maybe with cobra or I dont know but:

- one command switch like --use-qwen will have only certain subcommands and nothiung else
- one other sub command --use-flux to have other list of subcommands but not including the qwen ones etc. I need to orgnize it in such way

## CodeCompanion (Copilot)

Looking at your code, I can see you have three distinct model types that each need different flags:

- FLUX (default)
- Stable Diffusion (`--use-sd`)
- Qwen (`--use-qwen`)

For this use case, **Cobra** is perfect. Here's how to organize it with subcommands:

## Architecture Plan

```
gfluxgo
├── flux [default command]
│   ├── --config
│   ├── --gguf-model-url
│   ├── --hf-model
│   ├── --lora-url
│   ├── --img2img
│   └── ... (FLUX-specific flags)
├── sd
│   ├── --config
│   ├── --safetensors
│   ├── --sequential-offload
│   ├── --compile
│   └── ... (SD-specific flags)
└── qwen
    ├── --config
    ├── --qwen-num-images
    └── ... (Qwen-specific flags)
```

## Implementation

First, install Cobra:

```bash
go get -u github.com/spf13/cobra@latest
```

Create a new file structure:

```go {pkg/config/config.go}
// Package config - command line and all other configuration params
package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Seed          int
	OutputDir     string
	Resolution    []int
	Steps         int
	GuidanceScale float32
	CmdConf
}

type PromptConfig struct {
	StyleSuffix    string   `json:"style_suffix"`
	NegativePrompt string   `json:"negative_prompt"`
	Prompts        []string `json:"prompts"`
}

type CmdConf struct {
	ConfigPath        string
	ModelURL          string
	LoraURL           string
	ModelDownloadPath string
	LoraDownloadpath  string
	LowVRAM           bool
	HfModelID         string
	Debug             bool
	ImageToImage      bool
	Strength          float32
	WebServer         bool
	WebPort           int
	UseQwen           bool
	QwenNumImages     int
	SdCmd
}

type SdCmd struct {
	UseSD             bool
	SafetensorsPath   string
	SequentialOffload bool
	CompileModel      bool
	DisableSafety     bool
}

// Helper to get environment variables
func getEnv(key string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return ""
}

// ApplyEnvOverrides applies environment variable overrides to config
func (c *Config) ApplyEnvOverrides() error {
	if model := getEnv("MODEL_URL"); model != "" {
		c.ModelURL = model
	}

	if lora := getEnv("LORA_URL"); lora != "" {
		c.LoraURL = lora
	}

	if seed := getEnv("SEED"); seed != "" {
		s, err := strconv.ParseInt(seed, 10, 64)
		if err != nil {
			return fmt.Errorf("cannot parse SEED: %w", err)
		}
		c.Seed = int(s)
	}

	if hfModel := getEnv("HF_MODEL"); hfModel != "" {
		c.HfModelID = hfModel
	}

	if qwen := getEnv("USE_QWEN"); qwen != "" {
		if q, err := strconv.ParseBool(qwen); err == nil {
			c.UseQwen = q
		} else {
			log.Printf("Warning: cannot parse USE_QWEN from env: %v", err)
		}
	}

	if output := getEnv("OUTPUT"); output != "" {
		c.OutputDir = output
	}

	if st := getEnv("STEPS"); st != "" {
		steps, err := strconv.ParseInt(st, 10, 64)
		if err != nil {
			return fmt.Errorf("cannot parse STEPS: %w", err)
		}
		c.Steps = int(steps)
	}

	if str := getEnv("STRENGTH"); str != "" {
		if s, err := strconv.ParseFloat(str, 64); err == nil {
			c.Strength = float32(s)
		} else {
			return fmt.Errorf("cannot parse STRENGTH: %w", err)
		}
	}

	if guidance := getEnv("GUIDANCE_SCALE"); guidance != "" {
		if gs, err := strconv.ParseFloat(guidance, 64); err == nil {
			c.GuidanceScale = float32(gs)
		} else {
			return fmt.Errorf("cannot parse GUIDANCE_SCALE: %w", err)
		}
	}

	return nil
}

// ParseResolution parses resolution string like "1024x1024"
func ParseResolution(res string) ([]int, error) {
	r := strings.Split(res, "x")
	if len(r) != 2 {
		return nil, fmt.Errorf("invalid resolution format: %s, expected WIDTHxHEIGHT", res)
	}

	width, err := strconv.ParseInt(r[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse width: %w", err)
	}

	height, err := strconv.ParseInt(r[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse height: %w", err)
	}

	return []int{int(width), int(height)}, nil
}
```

Create the Cobra command structure:

```go {cmd/root.go}
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"yourmodule/pkg/config"
)

var cfg *config.Config

var rootCmd = &cobra.Command{
	Use:   "gfluxgo",
	Short: "AI image generation tool supporting FLUX, Stable Diffusion, and Qwen models",
	Long:  `A flexible image generation tool that supports multiple AI models including FLUX, Stable Diffusion, and Qwen.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	cfg = &config.Config{}

	// Global flags available to all commands
	rootCmd.PersistentFlags().StringVar(&cfg.ConfigPath, "config", "", "Path to JSON configuration file")
	rootCmd.PersistentFlags().StringVar(&cfg.OutputDir, "output", "./output", "Output directory for generated images")
	rootCmd.PersistentFlags().IntVar(&cfg.Seed, "seed", 42, "Random seed for generation")
	rootCmd.PersistentFlags().BoolVar(&cfg.Debug, "debug", false, "Enable debug mode")
	rootCmd.PersistentFlags().BoolVar(&cfg.WebServer, "web", false, "Enable web server mode")
	rootCmd.PersistentFlags().IntVar(&cfg.WebPort, "web-port", 8080, "Web server port")
	rootCmd.PersistentFlags().StringVar(&cfg.ModelDownloadPath, "model-path", "./models", "Download path for models")
}
```

```go {cmd/flux.go}
package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"yourmodule/pkg/config"
)

var (
	resolution    string
	guidanceScale float64
	strength      float64
)

var fluxCmd = &cobra.Command{
	Use:   "flux",
	Short: "Generate images using FLUX models",
	Long:  `Generate images using FLUX.1 models (dev or schnell variants)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Parse resolution
		res, err := config.ParseResolution(resolution)
		if err != nil {
			return err
		}
		cfg.Resolution = res
		cfg.GuidanceScale = float32(guidanceScale)
		cfg.Strength = float32(strength)

		// Apply environment overrides
		if err := cfg.ApplyEnvOverrides(); err != nil {
			return err
		}

		// Validation
		if cfg.ConfigPath == "" && !cfg.WebServer {
			return fmt.Errorf("config file required for image generation")
		}

		// Your FLUX generation logic here
		fmt.Printf("Running FLUX with config: %+v\n", cfg)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(fluxCmd)

	// FLUX-specific flags
	fluxCmd.Flags().StringVar(&cfg.ModelURL, "gguf-model-url", "", "URL to download FLUX GGUF model")
	fluxCmd.Flags().StringVar(&cfg.HfModelID, "hf-model", "black-forest-labs/FLUX.1-dev", "HuggingFace model ID")
	fluxCmd.Flags().StringVar(&cfg.LoraURL, "lora-url", "", "URL to download LoRA safetensors")
	fluxCmd.Flags().StringVar(&cfg.LoraDownloadpath, "lora-path", "./models/lora", "LoRA download path")
	fluxCmd.Flags().IntVar(&cfg.Steps, "steps", 28, "Number of inference steps")
	fluxCmd.Flags().StringVar(&resolution, "resolution", "1024x1024", "Image resolution (WIDTHxHEIGHT)")
	fluxCmd.Flags().Float64Var(&guidanceScale, "guidance-scale", 7.0, "Guidance scale for generation")
	fluxCmd.Flags().BoolVar(&cfg.LowVRAM, "low-vram", false, "Enable CPU offload for low VRAM")
	fluxCmd.Flags().BoolVar(&cfg.ImageToImage, "img2img", false, "Enable image-to-image mode")
	fluxCmd.Flags().Float64Var(&strength, "strength", 0.75, "Transformation strength for img2img (0.0-1.0)")
}
```

```go {cmd/sd.go}
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"yourmodule/pkg/config"
)

var sdCmd = &cobra.Command{
	Use:   "sd",
	Short: "Generate images using Stable Diffusion models",
	Long:  `Generate images using Stable Diffusion models (SD 1.5/2.1/SDXL/SD3)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg.UseSD = true

		// Parse resolution
		res, err := config.ParseResolution(resolution)
		if err != nil {
			return err
		}
		cfg.Resolution = res
		cfg.GuidanceScale = float32(guidanceScale)
		cfg.Strength = float32(strength)

		// Apply environment overrides
		if err := cfg.ApplyEnvOverrides(); err != nil {
			return err
		}

		// Validation
		if cfg.ConfigPath == "" && !cfg.WebServer {
			return fmt.Errorf("config file required for image generation")
		}

		// Your SD generation logic here
		fmt.Printf("Running Stable Diffusion with config: %+v\n", cfg)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(sdCmd)

	// SD-specific flags
	sdCmd.Flags().StringVar(&cfg.SafetensorsPath, "safetensors", "", "Path to safetensors model file")
	sdCmd.Flags().IntVar(&cfg.Steps, "steps", 50, "Number of inference steps")
	sdCmd.Flags().StringVar(&resolution, "resolution", "512x512", "Image resolution (WIDTHxHEIGHT)")
	sdCmd.Flags().Float64Var(&guidanceScale, "guidance-scale", 7.5, "Guidance scale for generation")
	sdCmd.Flags().BoolVar(&cfg.SequentialOffload, "sequential-offload", false, "Enable sequential CPU offload for ultra low VRAM")
	sdCmd.Flags().BoolVar(&cfg.CompileModel, "compile", false, "Compile model for faster inference")
	sdCmd.Flags().BoolVar(&cfg.DisableSafety, "disable-safety", true, "Disable safety checker")
	sdCmd.Flags().BoolVar(&cfg.ImageToImage, "img2img", false, "Enable image-to-image mode")
	sdCmd.Flags().Float64Var(&strength, "strength", 0.75, "Transformation strength for img2img (0.0-1.0)")
}
```

```go {cmd/qwen.go}
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"yourmodule/pkg/config"
)

var qwenCmd = &cobra.Command{
	Use:   "qwen",
	Short: "Edit images using Qwen-Image-Edit model",
	Long:  `Perform image editing tasks using the Qwen-Image-Edit model`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg.UseQwen = true

		// Parse resolution
		res, err := config.ParseResolution(resolution)
		if err != nil {
			return err
		}
		cfg.Resolution = res

		// Apply environment overrides
		if err := cfg.ApplyEnvOverrides(); err != nil {
			return err
		}

		// Validation
		if cfg.ConfigPath == "" {
			return fmt.Errorf("config file required for Qwen")
		}

		// Your Qwen generation logic here
		fmt.Printf("Running Qwen with config: %+v\n", cfg)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(qwenCmd)

	// Qwen-specific flags
	qwenCmd.Flags().IntVar(&cfg.QwenNumImages, "num-images", 1, "Number of images to generate per prompt")
	qwenCmd.Flags().StringVar(&resolution, "resolution", "1024x1024", "Image resolution (WIDTHxHEIGHT)")
	qwenCmd.Flags().IntVar(&cfg.Steps, "steps", 20, "Number of inference steps")
}
```

Update your main.go:

```go {main.go}
package main

import (
	"yourmodule/cmd"
)

func main() {
	cmd.Execute()
}
```

## Usage Examples

```bash
# FLUX command with its flags
./gfluxgo flux --config config.json --hf-model black-forest-labs/FLUX.1-dev --steps 28

# SD command with its flags
./gfluxgo sd --config config.json --safetensors model.safetensors --compile

# Qwen command with its flags
./gfluxgo qwen --config config.json --num-images 4

# Web server mode (works with all commands)
./gfluxgo flux --web --output ./output
```

This architecture ensures:

- Each model type has only its relevant flags
- No flag pollution between commands
- Clear separation of concerns
- Easy to extend with new models

Would you like me to show you how to add default command behavior or add flag validation?

## Me

> Context:
>
> -  <buf>pkg/config/config.go</buf>
