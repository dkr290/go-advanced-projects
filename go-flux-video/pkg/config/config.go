// Package config - comomand line and allo other configuation params
package config

import (
	"flag"
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
	ImageToImage      bool    // Enable img2img mode
	Strength          float32 // Transformation strength for img2img
	WebServer         bool    // Enable web server mode
	WebPort           int     // Web server port
	SdCmd
}

type SdCmd struct {
	UseSD             bool   // Use Stable Diffusion scripts instead of FLUX
	SafetensorsPath   string // Path to safetensors file (Civitai models)
	SequentialOffload bool   // Ultra low VRAM mode
	CompileModel      bool   // Torch compile for speed
	DisableSafety     bool   // Disable safety checker
}

func LoadConfig() *Config {
	c := &Config{}
	c.GetFlags()
	return c
}

func (c *Config) GetFlags() {
	var res string
	var guidanceScale float64
	flag.StringVar(
		&c.ConfigPath,
		"config",
		"",
		"Path to the JSON configuration file.",
	)

	flag.StringVar(
		&c.ModelURL,
		"gguf-model-url",
		"",
		"URL to download the FLUX GGUF model if not found locally.",
	)

	flag.StringVar(
		&c.HfModelID,
		"hf-model",
		"black-forest-labs/FLUX.1-dev",
		"HuggingFace model ID (default: FLUX.1-dev)",
	)

	flag.BoolVar(&c.Debug, "debug", false, "Using debug true or false")

	flag.StringVar(
		&c.LoraURL,
		"lora-url",
		"",
		"URL to download the LoRA safetensors file if not found locally.",
	)

	flag.IntVar(&c.Seed, "seed", 42, "seed number")
	flag.StringVar(
		&c.OutputDir,
		"output",
		"./output",
		"the output directory for the generated piuctures",
	)
	flag.StringVar(&res, "resolution", "1024x1024", "The resolution default 1024")

	flag.IntVar(
		&c.Steps,
		"num_steps",
		28,
		"number steps for the model, depends on the model",
	)

	flag.Float64Var(
		&guidanceScale,
		"guidence_scale",
		7.0,
		"number steps for the model, depends on the model",
	)

	flag.StringVar(
		&c.ModelDownloadPath,
		"model-down-path",
		"./models",
		"Download path of the models",
	)

	flag.StringVar(
		&c.LoraDownloadpath,
		"lora-down-path",
		"./models/lora",
		"Download path of the lora",
	)
	flag.BoolVar(
		&c.LowVRAM,
		"low-vram",
		false,
		"Enable CPU offload for low VRAM GPUs (slower on high VRAM GPUs)",
	)

	// Image-to-Image flags
	flag.BoolVar(
		&c.ImageToImage,
		"img2img",
		false,
		"Enable image-to-image mode instead of text-to-image",
	)

	var strength float64
	flag.Float64Var(
		&strength,
		"strength",
		0.75,
		"Transformation strength for img2img (0.0-1.0). Higher = more creative, lower = closer to input",
	)

	// Web server flags
	flag.BoolVar(
		&c.WebServer,
		"web",
		false,
		"Enable web server mode to view images in browser",
	)

	flag.IntVar(
		&c.WebPort,
		"web-port",
		8080,
		"Web server port (default: 8080)",
	)

	// flags fro SD models

	flag.BoolVar(
		&c.UseSD,
		"use-sd",
		false,
		"Use Stable Diffusion models instead of FLUX (SD 1.5/2.1/SDXL/SD3)",
	)

	flag.StringVar(
		&c.SafetensorsPath,
		"safetensors",
		"",
		"Path to safetensors model file (for Civitai models)",
	)

	flag.BoolVar(
		&c.SequentialOffload,
		"sequential-offload",
		false,
		"Enable sequential CPU offload for ultra low VRAM (4GB)",
	)
	flag.BoolVar(
		&c.CompileModel,
		"compile",
		false,
		"Compile model for 2x faster inference (requires PyTorch 2.0+)",
	)
	flag.BoolVar(
		&c.DisableSafety,
		"disable-safety",
		true,
		"Disable safety checker (allows NSFW content)",
	)

	flag.Parse()

	// Set strength after parsing
	c.Strength = float32(strength)

	if c.ConfigPath == "" {
		fmt.Println("Need the configuration file")
		os.Exit(1)
	}

	if model := getEnv("MODEL_URL"); model != "" {
		c.ModelURL = model
	}

	if lora := getEnv("LORA_URL"); lora != "" {
		c.LoraURL = lora
	}

	if seed := getEnv("SEED"); seed != "" {
		s, err := strconv.ParseInt(seed, 10, 64)
		if err != nil {
			log.Fatalf("cannot parse the seed %v", err)
		}
		c.Seed = int(s)

	}
	if hfModel := getEnv("HF_MODEL"); hfModel != "" {
		c.HfModelID = hfModel
	}

	if output := getEnv("OUTPUT"); output != "" {
		c.OutputDir = output
	}

	if st := getEnv("STEPS"); st != "" {
		steps, err := strconv.ParseInt(st, 10, 64)
		if err != nil {
			log.Fatalf("cannot parse the steps %v", err)
		}
		c.Steps = int(steps)

	}

	// Strength - environment variable takes precedence
	if str := getEnv("STRENGTH"); str != "" {
		if s, err := strconv.ParseFloat(str, 64); err == nil {
			c.Strength = float32(s)
		} else {
			log.Fatalf("cannot parse STRENGTH from env: %v", err)
		}
	}

	if resolution := getEnv("RESOLUTION"); resolution != "" {
		r := strings.Split(resolution, "x")
		if len(r) != 2 {
			log.Fatalf("invalid resolution format: %s, expected WIDTHxHEIGHT", res)
		}

		width, err := strconv.ParseInt(r[0], 10, 64)
		if err != nil {
			log.Fatalf("cannot parse the resolution %v", err)
		}

		height, err := strconv.ParseInt(r[1], 10, 64)
		if err != nil {
			log.Fatalf("cannot parse the resolution %v", err)
		}
		c.Resolution = []int{int(width), int(height)}

	} else {

		r := strings.Split(res, "x")
		if len(r) != 2 {
			log.Fatalf("invalid resolution format: %s, expected WIDTHxHEIGHT", res)
		}

		width, err := strconv.ParseInt(r[0], 10, 64)
		if err != nil {
			log.Fatalf("cannot parse the resolution %v", err)
		}

		height, err := strconv.ParseInt(r[1], 10, 64)
		if err != nil {
			log.Fatalf("cannot parse the resolution %v", err)
		}
		c.Resolution = []int{int(width), int(height)}

	}

	// Parse guidance scale - environment variable takes precedence
	if guidance := getEnv("GUIDANCE_SCALE"); guidance != "" {
		if gs, err := strconv.ParseFloat(guidance, 64); err == nil {
			c.GuidanceScale = float32(gs)
		} else {
			log.Fatalf("cannot parse GUIDANCE_SCALE from env: %v", err)
		}
	} else {
		c.GuidanceScale = float32(guidanceScale)
	}
}

func getEnv(key string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return ""
}
