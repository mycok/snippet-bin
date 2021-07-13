package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/mycok/snippet-bin/pkg/models"
)

func (app *application) home(wr http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(wr, err)

		return
	}

	app.render(wr, r, "home.page.go.tmpl", &templateData{Snippets: snippets})
}

func (app *application) showSnippet(wr http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id < 1 {
		app.notFoundError(wr)

		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNorRecord) {
			app.notFoundError(wr)
		} else {
			app.serverError(wr, err)
		}

		return
	}

	app.render(wr, r, "show.page.go.tmpl", &templateData{Snippet: snippet})
}

func (app *application) createSnippetForm(wr http.ResponseWriter, r *http.Request) {
	app.render(wr, r, "create.page.go.tmpl", nil)
}

func (app *application) createSnippet(wr http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		wr.Header().Set("Allow", http.MethodPost)
		app.clientError(wr, http.StatusMethodNotAllowed)
	
		return
	}

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := "7"
	
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(wr, err)

		return
	}

	// redirect the user to the relevant page
	http.Redirect(wr, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}


