package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// FileWriterResult represents the response structure of the tool
type FileWriterResult struct {
	Status string `json:"status"`
	File   string `json:"file"`
}
type FileWriteTool struct{}

func (FileWriteTool) Name() string {
	return "file_write"
}

func (FileWriteTool) Description() string {
	return "Writes content to a file and returns the path"
}

func (FileWriteTool) IsLongRunning() bool {
	return false
}

func (w *FileWriteTool) FileWriteTool(
	content string,
) (FileWriterResult, error) {
	root, err := os.Getwd()
	if err != nil {
		return FileWriterResult{Status: "error"}, err
	}

	outputDir := filepath.Join(root, "output")

	// Ensure the "output" directory exists.
	// 0755 is standard permissions (drwxr-xr-x)
	err = os.MkdirAll(outputDir, 0o755)
	if err != nil {
		return FileWriterResult{Status: "error"}, fmt.Errorf("failed to create directory: %w", err)
	}

	// Get current time formatted as YYMMDD_HHMMSS
	timestamp := time.Now().Format("060102_150405")

	// Construct the filename
	filename := fmt.Sprintf("%s_generated_page.html", timestamp)
	filePath := filepath.Join(outputDir, filename)

	// Write the content to the file.
	// 0644 gives read/write to owner and read-only to others.
	err = os.WriteFile(filePath, []byte(content), 0o644)
	if err != nil {
		return FileWriterResult{Status: "error"}, fmt.Errorf("failed to write file: %w", err)
	}

	return FileWriterResult{
		Status: "success",
		File:   filePath,
	}, nil
}
