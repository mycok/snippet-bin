package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(app.recoverFromPanic, app.logRequest, secureHeaders)

	mux.Route("/", func(mux chi.Router) {
		mux.With(app.session.Enable).Get("/", app.home)
		mux.With(app.session.Enable).Post("/snippet/create", app.createSnippet)
		mux.With(app.session.Enable).Get("/snippet/create", app.createSnippetForm)
		mux.With(app.session.Enable).Get("/snippet/{id}", app.showSnippet)
	})

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
