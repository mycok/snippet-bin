package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/mycok/snippet-bin/pkg/forms"
	"github.com/mycok/snippet-bin/pkg/models"
)

type templateData struct {
	CSRFToken string
	IsAuthenticated bool
	CurrentYear int
	Flash string
	Snippet *models.Snippet
	Snippets []*models.Snippet
	Form *forms.Form
}

// create a humanDate function which returns a nicely formatted string
// representation of a time.Time object.
func humanReadableDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}
// initialize a template.FuncMap object as a global variable
// it's used to map our custom template function string names with the actual functions
var functions = template.FuncMap{
	"humanReadableDate": humanReadableDate,
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
		// The template.FuncMap must be registered with the template set before you 
		// call the ParseFiles() method. This means we have to use template.New() to 
		// create an empty template set, use the Funcs() method to register the
		// template.FuncMap, and then parse the file as normal.
		templateSet, err := template.New(name).Funcs(functions).ParseFiles(page)
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
