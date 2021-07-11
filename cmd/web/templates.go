package main

import (
	"html/template"
	"path/filepath"

	"github.com/mycok/snippet-bin/pkg/models"
)

type templateData struct {
	CurrentYear int
	Snippet *models.Snippet
	Snippets []*models.Snippet
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	// Use the filepath.Glob function to get a slice of all filepaths with
	// the extension '.page.tmpl'. This essentially gives us a slice of all the
	// 'page' templates for the application.
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.go.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// Extract the file name (like 'home.page.tmpl') from the full file path
		// and assign it to the name variable.
		name := filepath.Base(page)
		// Parse the page template file in to a template set.
		templateSet, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Use the ParseGlob method to add any 'layout' templates to the 
		// template set (in our case, it's just the 'base' layout at the 
		// moment.
		templateSet, err = templateSet.ParseGlob(filepath.Join(dir, "*.layout.go.tmpl"))
		if err != nil {
			return nil, err
		}

		// Use the ParseGlob method to add any 'partial' templates to the
		// template set (in our case, it's just the 'footer' partial at the
		// moment.
		templateSet, err = templateSet.ParseGlob(filepath.Join(dir, "*.partial.go.tmpl"))
		if err != nil {
			return nil, err
		}

		// Add the template set to the cache, using the name of the page
		// (like 'home.page.go.tmpl') as the key.
		cache[name] = templateSet
	}

	return cache, nil
}
