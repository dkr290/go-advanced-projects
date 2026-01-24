package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

// FileWriteArgs defines the input parameters for the file write tool
type FileWriteArgs struct {
	Content string `json:"content" description:"The HTML content to write to a file"`
}

// FileWriteResult defines the output of the file write tool
type FileWriteResult struct {
	Status string `json:"status" description:"Status of the operation (success or error)"`
	File   string `json:"file"   description:"Path to the created file"`
}

// NewFileWriteTool creates a new file write tool using functiontool
func NewFileWriteTool() (tool.Tool, error) {
	handler := func(ctx tool.Context, input FileWriteArgs) (FileWriteResult, error) {
		fmt.Printf("[DEBUG] FileWriteTool called with content length: %d\n", len(input.Content))
		root, err := os.Getwd()
		if err != nil {
			fmt.Printf("[DEBUG] Failed to get working directory: %v\n", err)
			return FileWriteResult{
					Status: "error",
				}, fmt.Errorf(
					"failed to get working directory: %w",
					err,
				)
		}
		fmt.Printf("[DEBUG] Current working directory: %s\n", root)

		outputDir := filepath.Join(root, "output")
		fmt.Printf("[DEBUG] Output directory: %s\n", outputDir)
		err = os.MkdirAll(outputDir, 0o755)
		if err != nil {
			fmt.Printf("[DEBUG] Failed to create output directory: %v\n", err)
			return FileWriteResult{
					Status: "error",
				}, fmt.Errorf(
					"failed to create output directory: %w",
					err,
				)
		}

		timestamp := time.Now().Format("060102_150405")
		filename := fmt.Sprintf("%s_generated_page.html", timestamp)
		filePath := filepath.Join(outputDir, filename)

		fmt.Printf("[DEBUG] Writing file: %s\n", filePath)
		fmt.Printf("[DEBUG] Content preview (first 100 chars): %.100s\n", input.Content)

		err = os.WriteFile(filePath, []byte(input.Content), 0o644)
		if err != nil {
			fmt.Printf("[DEBUG] Failed to write file: %v\n", err)
			return FileWriteResult{Status: "error"}, fmt.Errorf("failed to write file: %w", err)
		}
		fmt.Printf("[DEBUG] File written successfully: %s\n", filePath)
		return FileWriteResult{
			Status: "success",
			File:   filePath,
		}, nil
	}

	return functiontool.New(functiontool.Config{
		Name:        "file_write",
		Description: "Writes HTML content to a file in the output directory and returns the file path",
	}, handler)
}
