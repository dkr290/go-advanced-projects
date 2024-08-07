package main

import (
	"dkr290/go-advanced-projects/snippets-toolbox/internal/models"
	"log"
)

type config struct {
	addr      string
	staticDir string
	dsn       string
}

// adding snippets fields into the application struct this way it will allow to make SnippetModel object
// available to our handlers
type appconfig struct {
	errotLog *log.Logger
	infoLog  *log.Logger
	snippets *models.SnippetsModel
}
