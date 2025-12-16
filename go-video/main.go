package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/binozo/gostablediffusion/pkg/sd"
)

type Config struct {
	Seed           int      `json:"seed"`
	OutputDir      string   `json:"output_dir"`
	Resolution     []int    `json:"resolution"`
	Steps          int      `json:"num_inference_steps"`
	GuidanceScale  float32  `json:"guidance_scale"`
	StyleSuffix    string   `json:"style_suffix"`
	NegativePrompt string   `json:"negative_prompt"`
	Prompts        []string `json:"prompts"`
}

type WriteCounter struct {
	Total int64 // Total bytes written
	Size  int64 // Total file size
	last  time.Time
}

// Write implements the io.Writer interface.
func (wc *WriteCounter) Write(p []byte) (n int, err error) {
	n = len(p)
	wc.Total += int64(n)
	now := time.Now()
	if now.Sub(wc.last) > 500*time.Millisecond {
		wc.PrintProgress()
		wc.last = now
	}
	return
}

// PrintProgress displays the download percentage
func (wc *WriteCounter) PrintProgress() {
	if wc.Size > 0 {
		percent := float64(wc.Total) / float64(wc.Size) * 100
		// \r returns the cursor to the start of the line, allowing overwrite
		fmt.Printf("\rDownloading... %.2f%% (%s / %s)",
			percent,
			byteCountToHuman(wc.Total),
			byteCountToHuman(wc.Size))
	} else {
		// Fallback if Content-Length is missing
		fmt.Printf("\rDownloading... %s (size unknown)", byteCountToHuman(wc.Total))
	}
}

func byteCountToHuman(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

// DownloadFile handles the download, saving it to a local file.
func DownloadFile(url, destPath string) error {
	fmt.Printf("Attempting to download from: %s\n", url)

	// 1. HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad HTTP status: %s", resp.Status)
	}

	// 2. Setup the output file
	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", destPath, err)
	}
	defer out.Close()

	// 3. Write and track progress
	counter := &WriteCounter{Size: resp.ContentLength}

	// TeeReader pipes the data through the counter while copying to the file
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		return fmt.Errorf("failed to write file content: %w", err)
	}

	// Final progress update to ensure 100% is displayed
	counter.Total = counter.Size // Force total to size for final display
	counter.PrintProgress()
	fmt.Println("\nDownload complete.")
	return nil
}

// --- Main Application Logic ---

func sanitizeFilename(text string) string {
	text = strings.ReplaceAll(text, "/", "-")
	text = strings.ReplaceAll(text, " ", "-")
	text = strings.ReplaceAll(text, ",", "")
	return strings.ToLower(text)
}

func main() {
	// --- COMMAND LINE FLAG DEFINITION ---
	var configPath string
	var modelURL string
	var loraURL string

	flag.StringVar(
		&configPath,
		"config",
		"character_config.json",
		"Path to the JSON configuration file.",
	)
	flag.StringVar(
		&modelURL,
		"model-url",
		"https://huggingface.co/city96/FLUX.1-dev-gguf/resolve/main/flux1-dev-Q4_K_S.gguf",
		"URL to download the FLUX GGUF model if not found locally.",
	)
	flag.StringVar(
		&loraURL,
		"lora-url",
		"https://huggingface.co/Heartsync/Flux-NSFW-uncensored/resolve/main/lora.safetensors",
		"URL to download the LoRA safetensors file if not found locally.",
	)
	flag.Parse()

	// --- Configuration Loading ---
	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("Error reading config file '%s': %v\n", configPath, err)
		os.Exit(1)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	// --- Model Paths (Local) ---
	// The models are stored locally relative to the execution directory
	modelPath := filepath.Join("./models", getFilenameFromURL(modelURL))
	loraDir := "./models/lora"
	loraPath := filepath.Join("./models/lora", getFilenameFromURL(loraURL))

	// Ensure the 'lib' directory exists
	if err := os.MkdirAll(filepath.Dir(modelPath), 0o755); err != nil {
		fmt.Printf("Error creating model directory: %v\n", err)
		os.Exit(1)
	}
	if err := os.MkdirAll(filepath.Dir(loraPath), 0o755); err != nil {
		fmt.Printf("Error creating lora directory: %v\n", err)
		os.Exit(1)
	}

	// --- 1. CONDITIONAL DOWNLOAD LOGIC (Model) ---
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		if modelURL == "" {
			fmt.Printf(
				"FATAL: Model file not found at '%s', and no --model-url was provided.\n",
				modelPath,
			)
			os.Exit(1)
		}
		fmt.Printf("Model not found locally. Downloading to '%s'...\n", modelPath)
		if err := DownloadFile(modelURL, modelPath); err != nil {
			fmt.Printf("FATAL Download Error: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Model found locally at: %s\n", modelPath)
	}
	// --- 3. CONDITIONAL DOWNLOAD & LOAD LOGIC (LoRA) ---
	if _, err := os.Stat(loraPath); os.IsNotExist(err) {
		if loraURL != "" {
			fmt.Printf("LoRA file not found locally. Downloading to '%s'...\n", loraPath)
			if err := DownloadFile(loraURL, loraPath); err != nil {
				fmt.Printf("LoRA Download Warning: %v. Continuing without LoRA.\n", err)
			}
		} else {
			fmt.Println("LoRA file not found, and no --lora-url provided. Skipping LoRA.")
		}
	}

	// --- 2. Initialize the Model (Calls C++ code via CGO) ---
	fmt.Println("Starting FLUX model initialization via CGO...")

	sdBuilder := sd.New().SetModel(modelPath)
	sdBuilder.SetLoRaDir(loraDir)
	ctx, err := sdBuilder.Load()
	if err != nil {
		fmt.Printf("FATAL CGO Model Load Error (Is your GGUF file correct?): %v\n", err)
		os.Exit(1)
	}
	defer ctx.Free() // Releases C++ resources

	// --- 4. Generation Loop ---
	outputDir := cfg.OutputDir
	os.MkdirAll(outputDir, 0o755)

	fmt.Printf("\nStarting character sheet generation for %d poses...\n", len(cfg.Prompts))

	for i, p := range cfg.Prompts {
		// Construct the final prompt
		prompt := fmt.Sprintf("%s, %s", p, cfg.StyleSuffix)

		fmt.Printf("[%d/%d] Generating: %s\n", i+1, len(cfg.Prompts), p)

		// Create the generation parameters
		params := sd.NewImageGenerationParams()
		params.Width = cfg.Resolution[0]
		params.Height = cfg.Resolution[1]
		params.Prompt = prompt
		params.NegativePrompt = cfg.NegativePrompt
		params.SampleSteps = cfg.Steps
		params.Guidance.TxtCfg = cfg.GuidanceScale
		params.Seed = int64(cfg.Seed + i)

		start := time.Now()
		// CGO CALL: The generation happens here
		imageData := ctx.GenerateImage(params)
		// Save Image
		safeLabel := sanitizeFilename(p)
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

	fmt.Println("\n✅ Go Character Sheet Generation Complete!")
}

func getFilenameFromURL(url string) string {
	// Extract the last part of the URL
	parts := strings.Split(url, "/")
	return parts[len(parts)-1]
}
