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
	Debug             bool
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
		"model-url",
		"https://huggingface.co/city96/FLUX.1-dev-gguf/resolve/main/flux1-dev-Q4_K_S.gguf",
		"URL to download the FLUX GGUF model if not found locally.",
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
		"/lib/output",
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
	flag.Parse()
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
