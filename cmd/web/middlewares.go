package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

// To be executed on every request
func secureHeaders(next http.Handler)  http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		wr.Header().Set("x-XSS-Protection", "1; mode=block")
		wr.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(wr, r)
	})
}

func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path: "/",
		Secure: true,
	})

	return csrfHandler

}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(wr, r)
	})
}

func (app *application) recoverFromPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event
		// of a panic as Go unwinds the stack).
		defer func() {
			// Use the builtin recover function to check if there has been a 
			// panic or not. If there has...

			if err := recover(); err != nil {
				wr.Header().Set("Connection", "Close")
				app.serverError(wr, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(wr, r)
	})
}

func (app *application) requireAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(rw, r, "/login", http.StatusSeeOther)

			return
		}
		// set the "Cache-Control: no-store" header so that pages that
		// require authentication are not stored in the users browser cache 
		// other intermediary cache
		rw.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(rw, r)
	})
}
