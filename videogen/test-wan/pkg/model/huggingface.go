package model

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/wan2-video-server/pkg/config"
	"github.com/wan2-video-server/pkg/logger"
)

// HuggingFaceDownloader handles downloading models from Hugging Face
type HuggingFaceDownloader struct {
	config *config.Config
	log    *logger.Logger
}

// NewHuggingFaceDownloader creates a new Hugging Face downloader
func NewHuggingFaceDownloader(cfg *config.Config) *HuggingFaceDownloader {
	return &HuggingFaceDownloader{
		config: cfg,
		log:    logger.NewLogger(),
	}
}

// Download downloads the model from Hugging Face
func (d *HuggingFaceDownloader) Download() error {
	d.log.Infof("Downloading model from Hugging Face: %s", d.config.Model.HuggingFaceModelID)

	// Create cache directory
	cacheDir := filepath.Join(d.config.Model.CacheDir, d.config.Model.HuggingFaceModelID)
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Files to download (example for LTX-Video model)
	files := []string{
		"config.json",
		"model.safetensors",
		"tokenizer_config.json",
		"tokenizer.json",
		"vocab.json",
		"merges.txt",
	}

	baseURL := fmt.Sprintf("%s/%s/resolve/main", d.config.HF.APIURL, d.config.Model.HuggingFaceModelID)

	for _, file := range files {
		if err := d.downloadFile(baseURL, file, cacheDir); err != nil {
			d.log.Warnf("Failed to download %s: %v (might not exist)", file, err)
			// Continue with other files
		}
	}

	d.log.Info("Model download completed")
	return nil
}

func (d *HuggingFaceDownloader) downloadFile(baseURL, filename, destDir string) error {
	url := fmt.Sprintf("%s/%s", baseURL, filename)
	destPath := filepath.Join(destDir, filename)

	// Check if file already exists
	if _, err := os.Stat(destPath); err == nil {
		d.log.Infof("File already exists, skipping: %s", filename)
		return nil
	}

	d.log.Infof("Downloading: %s", filename)

	// Create HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// Add authorization header if token is provided
	if d.config.HF.Token != "" && d.config.HF.Token != "your_hf_token_here" {
		req.Header.Set("Authorization", "Bearer "+d.config.HF.Token)
	}

	// Download file
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download: status %d", resp.StatusCode)
	}

	// Create destination file
	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Copy data
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	d.log.Infof("Downloaded: %s", filename)
	return nil
}
