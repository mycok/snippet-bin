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

func (app *application) home(rw http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(rw, err)

		return
	}

	app.render(rw, r, "home.page.go.tmpl", &templateData{Snippets: snippets})
}

func (app *application) showSnippet(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id < 1 {
		app.notFoundError(rw)

		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNorRecord) {
			app.notFoundError(rw)
		} else {
			app.serverError(rw, err)
		}

		return
	}

	app.render(rw, r, "show.page.go.tmpl", &templateData{Snippet: snippet})
}

func (app *application) createSnippetForm(rw http.ResponseWriter, r *http.Request) {
	app.render(rw, r, "create.page.go.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) createSnippet(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(rw, http.StatusBadRequest)

		return
	}

	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	if !form.IsValid() {
		app.render(rw, r, "create.page.go.tmpl", &templateData{Form: form})

		return
	}
	
	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(rw, err)

		return
	}

	app.session.Put(r, "flash", "Snippet successfully created!")
	// redirect the user to the relevant page
	http.Redirect(rw, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) signupForm(rw http.ResponseWriter, r *http.Request) {
	app.render(rw, r, "signup.page.go.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) signup(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(rw, http.StatusBadRequest)

		return
	}

	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MaxLength("name", 255)
	form.MaxLength("email", 255)
	form.MatchesPattern("email", forms.EmailRegex)
	form.MinLength("password", 10)

	if !form.IsValid() {
		app.render(rw, r, "signup.page.go.tmpl", &templateData{
			Form: form,
		})

		return
	}

	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors.Add("email", "Email already in use")
			app.render(rw, r, "signup.page.go.tmpl", &templateData{
				Form: form,
			})
		} else {
			app.serverError(rw, err)
		}

		return
	}

	app.session.Put(r, "flash", "Your signup was successful. please log in.")

	http.Redirect(rw, r, "/login", http.StatusSeeOther)

}

func (app *application) loginForm(rw http.ResponseWriter, r *http.Request) {
	app.render(rw, r, "login.page.go.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) login(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(rw, http.StatusBadRequest)
	}

	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or Password is incorrect")
			app.render(rw, r, "login.page.go.tmpl", &templateData{
				Form: form,
			})
		} else {
			app.serverError(rw, err)
		}

		return
	}

	app.session.Put(r, "authenticatedUserID", id)

	http.Redirect(rw, r, "/snippet/create", http.StatusSeeOther)
}

func (app *application) logout(rw http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "authenticatedUserID")
	app.session.Put(r, "flash", "You've been successfully logged out")
	
	http.Redirect(rw, r, "/", http.StatusSeeOther)
}
