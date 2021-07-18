package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/mycok/snippet-bin/pkg/forms"
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
	app.render(wr, r, "create.page.go.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) createSnippet(wr http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(wr, http.StatusBadRequest)

		return
	}

	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	if !form.IsValid() {
		app.render(wr, r, "create.page.go.tmpl", &templateData{Form: form})

		return
	}
	
	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(wr, err)

		return
	}

	app.session.Put(r, "flash", "Snippet successfully created!")
	// redirect the user to the relevant page
	http.Redirect(wr, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) signupForm(wr http.ResponseWriter, r *http.Request) {
	app.render(wr, r, "signup.page.go.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) signup(wr http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(wr, http.StatusBadRequest)

		return
	}

	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MaxLength("name", 255)
	form.MaxLength("email", 255)
	form.MatchesPattern("email", forms.EmailRegex)
	form.MinLength("password", 10)

	if !form.IsValid() {
		app.render(wr, r, "signup.page.go.tmpl", &templateData{
			Form: form,
		})

		return
	}

	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors.Add("email", "Email already in use")
			app.render(wr, r, "signup.page.go.tmpl", &templateData{
				Form: form,
			})
		} else {
			app.serverError(wr, err)
		}

		return
	}

	app.session.Put(r, "flash", "Your signup was successful. please log in.")

	http.Redirect(wr, r, "/login", http.StatusSeeOther)

}

func (app *application) loginForm(wr http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(wr, "Display user login form")
}

func (app *application) login(wr http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(wr, "Login user")
}

func (app *application) logout(wr http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(wr, "Logout user")
}



