package main

import (
	"html/template"
	"path/filepath"

	"poc.eduardo-luz.eu/internal/models"
)

// acts as a holding structure for dynamic data that will be available to the templates
type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// gives a slice of all the filepaths for our application 'page' templates
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		//extract the filename from the fullpath
		name := filepath.Base(page)
		//slice with the filepaths for the base template and any partials
		files := []string{
			"./ui/html/base.tmpl",
			"./ui/html/partials/nav.tmpl",
			page,
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil

}
