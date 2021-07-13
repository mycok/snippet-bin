package main

import (
	"fmt"
	"net/http"
)

// To be executed on every request
func secureHeaders(next http.Handler)  http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		wr.Header().Set("x-XSS-Protection", "1; mode=block")
		wr.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(wr, r)
	})
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