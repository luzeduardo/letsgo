package main

import (
	"html/template"
	"path/filepath"
	"time"

	"poc.eduardo-luz.eu/internal/models"
)

// acts as a holding structure for dynamic data that will be available to the templates
type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// map of functions that we want to make available to all our templates
// template functions can return only 1 value
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// gives a slice of all the filepaths for our application 'page' templates
	pages, err := filepath.Glob("./ui/html/pages/*.go.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		//extract the filename from the fullpath
		name := filepath.Base(page)

		// FuncMap must be registered before calling ParseFiles
		// create a new template set, add the template functions and then parse files
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.go.tmpl")
		if err != nil {
			return nil, err
		}

		//ParseGlob to add any partials
		ts, err = ts.ParseGlob("./ui/html/partials/*.go.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
