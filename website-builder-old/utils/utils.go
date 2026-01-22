package utils

import (
	"fmt"
	"os"
)

// LoadInstructionsFile reads text from a file.
// If the file cannot be read, it logs a message and returns the provided default string.
func LoadInstructionsFile(filename string) (string, error) {
	// os.ReadFile is the idiomatic way to read an entire file into memory in Go.
	// It handles opening and closing the file automatically.
	data, err := os.ReadFile(filename)
	if err != nil {
		// If the error is specifically that the file doesn't exist
		if os.IsNotExist(err) {
			return "", fmt.Errorf("[WARNING] File not found: %s. Using default", filename)
		} else {
			// Catch any other errors (permissions, etc.)
			return "", fmt.Errorf("[ERROR] Failed to load %s: %v", filename, err)
		}
	}

	// os.ReadFile returns []byte, so we convert it to a string.
	// Go uses UTF-8 by default for string conversions.
	return string(data), nil
}
