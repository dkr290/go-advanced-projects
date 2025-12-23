package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"gfluxgo/pkg/config"
	"gfluxgo/pkg/generate"
	"gfluxgo/pkg/logging"
	"gfluxgo/pkg/utils"
	"gfluxgo/pkg/webserver"
)

func main() {
	var cmdConf config.Config
	var promptConf config.PromptConfig

	cmdConf.GetFlags()
	llogger := logging.Init(cmdConf.Debug)

	// --- Configuration Loading ---
	data, err := os.ReadFile(cmdConf.ConfigPath)
	if err != nil {
		fmt.Printf("Error reading config file '%s': %v\n", cmdConf.ConfigPath, err)
		os.Exit(1)
	}

	if err := json.Unmarshal(data, &promptConf); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	// --- Model Paths (Local) ---
	// The models are stored locally relative to the execution directory
	modelPath := filepath.Join(
		cmdConf.ModelDownloadPath,
		utils.GetFilenameFromURL(cmdConf.ModelURL),
	)
	loraDir := cmdConf.LoraDownloadpath
	loraPath := filepath.Join(loraDir, utils.GetFilenameFromURL(cmdConf.LoraURL))

	// Ensure the 'lib' directory exists
	if err := os.MkdirAll(filepath.Dir(modelPath), 0o755); err != nil {
		fmt.Printf("Error creating model directory: %v\n", err)
		os.Exit(1)
	}
	if err := os.MkdirAll(filepath.Dir(loraPath), 0o755); err != nil {
		fmt.Printf("Error creating lora directory: %v\n", err)
		os.Exit(1)
	}
	if err := utils.DownloadFiles(modelPath, cmdConf.ModelURL, cmdConf.LoraURL, loraPath, *llogger); err != nil {
		llogger.Logging.Errorf("error %v", err)
		os.Exit(1)
	}

	if cmdConf.ImageToImage {
		llogger.Logging.Infof("Starting Image to image mode ")
		inputImagesDir := "./images"

		// Choose SD or FLUX script
		if cmdConf.UseSD {
			if err := generate.GenerateImg2ImgWithPythonSD(cmdConf, promptConf, modelPath, loraDir, inputImagesDir); err != nil {
				llogger.Logging.Errorf("SD image to image generation failed: %v", err)
				os.Exit(1)
			}
		} else {
			if err := generate.GenerateImg2ImgWithPython(cmdConf, promptConf, modelPath, loraDir, inputImagesDir); err != nil {
				llogger.Logging.Errorf("Image to image generation failed: %v", err)
				os.Exit(1)
			}
		}
		fmt.Println("\n✅ Image-to-Image Generation Complete!")

	} else {
		// Choose SD or FLUX script
		if cmdConf.UseSD {
			llogger.Logging.Infof("Starting Stable Diffusion model initialization")
			if err := generate.GenerateWithPythonSD(cmdConf, promptConf, modelPath, loraDir); err != nil {
				llogger.Logging.Errorf("SD generate images failed %v", err)
				os.Exit(1)
			}
		} else {
			llogger.Logging.Infof(
				"Starting %s model initialization",
				utils.GetFilenameFromURL(cmdConf.ModelURL),
			)
			if err := generate.GenerateWithPython(cmdConf, promptConf, modelPath, loraDir); err != nil {
				llogger.Logging.Errorf("Generate images failed %v", err)
				os.Exit(1)
			}
		}

		fmt.Println("\n✅ Go image Generation Complete!")
	}
	// Start web server if enabled
	if cmdConf.WebServer {
		fmt.Println("\n🚀 Starting web server...")
		server := webserver.NewServer(cmdConf.OutputDir, cmdConf.WebPort)

		// Handle graceful shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		go func() {
			<-sigChan
			fmt.Println("\n\n👋 Shutting down web server...")
			os.Exit(0)
		}()

		if err := server.Start(); err != nil {
			llogger.Logging.Errorf("Web server failed: %v", err)
			os.Exit(1)
		}
	}
}
