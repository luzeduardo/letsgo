package main

import "poc.eduardo-luz.eu/internal/models"

// acts as a holding structure for dynamic data that will be available to the templates
type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}
