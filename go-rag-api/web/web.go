// Package web - using web interface for the app
package web

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"strings"

	"github.com/dkr290/go-advanced-projects/go-rag-api/llm"
	"github.com/dkr290/go-advanced-projects/go-rag-api/rag"
	"github.com/dkr290/go-advanced-projects/go-rag-api/vector"
	"github.com/gofiber/fiber/v3"
)

//go:embed templates/*.html
var templateFS embed.FS

type Options struct {
	Addr             string
	SystemPromptFile string
	Title            string
	Store            vector.Store
	ProcessDir       string
	ImagesDir        string
}
type Server struct {
	app            *fiber.App
	client       *llm.Client
	retreiver    *rag.Retriever
	processedDir string
	imagesDir    string
	tpl          *template.Template
	store vector.Store
	system       string
	title        string
}

func New(client *llm.Client, retreiver *rag.Retriever, options Options) (*Server, error) {
	tpl, err := template.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("parse templates error %w", err)
	}

	title := options.Title
	if title == "" {
		title = "Rag Chat"
	}

	return &Server{
		client:       client,
		title:        title,
		retreiver:    retreiver,
		processedDir: options.ProcessDir,
		imagesDir:    options.ImagesDir,
		store: options.Store,
		tpl: tpl,
		system: readSystemPrompt(options.SystemPromptFile),
	}, nil
}

func readSystemPrompt(path string) string {
  if path == "" {
		return ""
	}

	data,err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return ""
	}

	return strings.TrimSpace(string(data))
}

