package main

import (
	"encoding/json"
	"fmt"
	"gfluxgo/pkg/config"
	"gfluxgo/pkg/generate"
	"gfluxgo/pkg/logging"
	"gfluxgo/pkg/utils"
	"os"
	"path/filepath"
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

	llogger.Logging.Infof(
		"Starting %s model initialization",
		utils.GetFilenameFromURL(cmdConf.ModelURL),
	)

	if err := generate.GenerateWithPython(cmdConf, promptConf, modelPath, loraDir); err != nil {
		llogger.Logging.Errorf("Generate images failed %v", err)
		os.Exit(1)
	}

	fmt.Println("\n✅ Go image Generation Complete!")
}
