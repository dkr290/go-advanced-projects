package main

import "dkr290/go-advanced-projects/snippets-toolbox/internal/models"

//define templateData type to act as holding ctructure for any dynamic data
//to pass to the html templates

type TemplateData struct {
	Snippet *models.Snippet
}
