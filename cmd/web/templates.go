package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"poc.eduardo-luz.eu/internal/models"
	"poc.eduardo-luz.eu/ui"
)

// acts as a holding structure for dynamic data that will be available to the templates
type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

// map of functions that we want to make available to all our templates
// template functions can return only 1 value
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// gives a slice of all the filepaths for our application 'page' templates
	pages, err := fs.Glob(ui.Files, "html/pages/*.go.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		//extract the filename from the fullpath
		name := filepath.Base(page)

		// (not using File system embedded) FuncMap must be registered before calling ParseFiles
		// create a new template set, add the template functions and then parse files

		patterns := []string{
			"html/base.go.tmpl",
			"html/partials/*.go.tmpl",
			page,
		}

		// now it parse html directly from the Filesystem instead of the files (using embedded FS)
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
