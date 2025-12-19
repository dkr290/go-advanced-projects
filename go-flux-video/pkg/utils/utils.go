// Package utils
package utils

import (
	"fmt"
	"gfluxgo/pkg/download"
	"gfluxgo/pkg/logging"
	"os"
	"strings"
)

func SanitizeFilename(text string) string {
	text = strings.ReplaceAll(text, "/", "-")
	text = strings.ReplaceAll(text, " ", "-")
	text = strings.ReplaceAll(text, ",", "")
	return strings.ToLower(text)
}

func SanitizeFilenameForImage(prompt string, index int) string {
	// Sanitize the prompt
	sanitized := SanitizeFilename(prompt)

	// Take first 10 characters (or less) for the prompt part
	promptPart := sanitized
	if len(promptPart) > 10 {
		promptPart = promptPart[:10]
	}

	// Remove trailing dash if present
	promptPart = strings.TrimSuffix(promptPart, "-")

	// Format: {index}_{prompt10}.png
	// Max: 2 + 1 + 10 + 4 = 17 chars (well under 20)

	return fmt.Sprintf("%02d_%s.png", index, promptPart)
}

func GetFilenameFromURL(url string) string {
	// Extract the last part of the URL
	parts := strings.Split(url, "/")
	return parts[len(parts)-1]
}

func DownloadFiles(modelPath, modelURL, loraURL, loraPath string, l logging.Logger) error {
	// --- 1. CONDITIONAL DOWNLOAD LOGIC (Model) ---
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		if modelURL != "" {

			l.Logging.Infof("Model not found locally. Downloading to '%s'...\n", modelPath)
			if err := download.DownloadFile(modelURL, modelPath, l); err != nil {
				l.Logging.Errorf("FATAL Download Error: %v\n", err)
				return err
			}
		}
	} else {
		l.Logging.Infof("Model found locally at: %s\n", modelPath)
	}
	// --- 3. CONDITIONAL DOWNLOAD & LOAD LOGIC (LoRA) ---
	if _, err := os.Stat(loraPath); os.IsNotExist(err) {
		if loraURL != "" {
			l.Logging.Infof("LoRA file not found locally. Downloading to '%s'...\n", loraPath)
			if err := download.DownloadFile(loraURL, loraPath, l); err != nil {
				l.Logging.Errorf("LoRA Download Error: %v.\n", err)
				return err
			}
		} else {
			l.Logging.Warning("LoRA file not found, and no --lora-url provided. Skipping LoRA.")
		}
	}
	return nil
}
