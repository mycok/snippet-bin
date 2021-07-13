package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(app.recoverFromPanic, app.logRequest, secureHeaders)

	mux.Get("/", app.home)
	mux.Post("/snippet/create", app.createSnippet)
	mux.Get("/snippet/create", app.createSnippetForm)
	mux.Get("/snippet/{id}", app.showSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
