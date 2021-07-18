package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(app.recoverFromPanic, app.logRequest, secureHeaders)

	mux.Route("/", func(mux chi.Router) {
		mux.With(app.session.Enable, noSurf).Get("/", app.home)
		mux.With(app.session.Enable, noSurf).Get("/signup", app.signupForm)
		mux.With(app.session.Enable, noSurf).Post("/signup", app.signup)
		mux.With(app.session.Enable, noSurf).Get("/login", app.loginForm)
		mux.With(app.session.Enable, noSurf).Post("/login", app.login)
		mux.With(app.session.Enable, noSurf, app.requireAuthorization).Post("/logout", app.logout)
		mux.With(app.session.Enable, noSurf, app.requireAuthorization).Post("/snippet/create", app.createSnippet)
		mux.With(app.session.Enable, noSurf, app.requireAuthorization).Get("/snippet/create", app.createSnippetForm)
		mux.With(app.session.Enable, noSurf).Get("/snippet/{id}", app.showSnippet)
	})

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
