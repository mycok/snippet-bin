package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mycok/snippet-bin/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFoundError(w)

		return
	}

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)

		return
	}

	app.render(w, r, "home.page.go.tmpl", &templateData{Snippets: snippets})
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFoundError(w)

		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNorRecord) {
			app.notFoundError(w)
		} else {
			app.serverError(w, err)
		}

		return
	}

	app.render(w, r, "show.page.go.tmpl", &templateData{Snippet: snippet})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
	
		return
	}

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := "7"
	
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)

		return
	}

	// redirect the user to the relevant page
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}


