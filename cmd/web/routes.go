package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(app.recoverFromPanic, app.logRequest, secureHeaders)

	mux.Route("/", func(mux chi.Router) {
		mux.With(app.session.Enable, noSurf, app.authenticate).Get("/", app.home)
		mux.With(app.session.Enable, noSurf, app.authenticate).Get("/signup", app.signupForm)
		mux.With(app.session.Enable, noSurf, app.authenticate).Post("/signup", app.signup)
		mux.With(app.session.Enable, noSurf, app.authenticate).Get("/login", app.loginForm)
		mux.With(app.session.Enable, noSurf, app.authenticate).Post("/login", app.login)
		mux.With(app.session.Enable, noSurf, app.authenticate, app.requireAuthentication).Post("/logout", app.logout)
		mux.With(app.session.Enable, noSurf, app.authenticate, app.requireAuthentication).Post("/snippet/create", app.createSnippet)
		mux.With(app.session.Enable, noSurf, app.authenticate, app.requireAuthentication).Get("/snippet/create", app.createSnippetForm)
		mux.With(app.session.Enable, noSurf, app.authenticate).Get("/snippet/{id}", app.showSnippet)
	})

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
