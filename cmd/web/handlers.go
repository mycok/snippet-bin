package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

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
	err := r.ParseForm()
	if err != nil {
		app.clientError(wr, http.StatusBadRequest)

		return
	}	

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")

	errors := make(map[string]string)

	if strings.TrimSpace(title) == "" {
		errors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		errors["title"] = "This value is too long. (maximum allowed is 100 characters)"
	}

	if strings.TrimSpace(content) == "" {
		errors["content"] = "This field cannot be blank"
	}

	if strings.TrimSpace(expires) == "" {
		errors["expires"] = "This field cannot be blank"
	} else if expires != "365" && expires != "7" && expires != "1" {
		errors["expires"] = "This field is invalid"
	}

	if len(errors) > 0 {
		app.render(wr, r, "create.page.go.tmpl", &templateData{
			FormData: r.PostForm,
			FormErrors: errors,
		})

		return
	}
	
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(wr, err)

		return
	}

	// redirect the user to the relevant page
	http.Redirect(wr, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}


