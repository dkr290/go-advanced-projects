// Package template to generate templates from the yaml file
package template

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"crd-api-deploy/internal/models"
)

// Engine handles template generation
type Engine struct {
	template *template.Template
}

// NewEngine creates a new template engine
func NewEngine() (*Engine, error) {
	return NewEngineFromFile("templates/crd.yaml")
}

// NewEngineFromFile creates a new template engine loading from file path
func NewEngineFromFile(templatePath string) (*Engine, error) {
	if templatePath == "" {
		templatePath = "templates/crd.yaml"
	}

	templatePath = filepath.Clean(templatePath)

	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("template file not found: %s", templatePath)
	}

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template file %s: %w", templatePath, err)
	}

	return &Engine{
		template: tmpl,
	}, nil
}

// GenerateCRD generates a CRD YAML from the request
func (e *Engine) GenerateCRD(req *models.CreateAPIRequest) (string, error) {
	var buf bytes.Buffer

	// Set default values if not provided
	if req.Resources.Limits.EphemeralStorage == "" {
		req.Resources.Limits.EphemeralStorage = "6Gi"
	}

	// Execute template
	err := e.template.Execute(&buf, req)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
