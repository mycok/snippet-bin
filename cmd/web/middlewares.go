package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/mycok/snippet-bin/pkg/models"

	"github.com/justinas/nosurf"
)

// To be executed on every request
func secureHeaders(next http.Handler) http.Handler {
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
		Path:     "/",
		Secure:   true,
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
			// panic or not. If there was...

			if err := recover(); err != nil {
				wr.Header().Set("Connection", "Close")
				app.serverError(wr, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(wr, r)
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(rw, r, "/login", http.StatusSeeOther)

			return
		}
		// set the "Cache-Control: no-store" header so that pages that
		// require authentication are not stored in the users browser cache
		// or other intermediary cache
		rw.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(rw, r)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// Check if a authenticatedUserID value exists in the session. If this *isn't
		// present* then call the next handler in the chain as normal.
		exists := app.session.Exists(r, "authenticatedUserID")
		if !exists {
			next.ServeHTTP(rw, r)

			return
		}
		// Fetch the details of the current user from the database. If no matching
		// record is found, or the current user is has been deactivated, remove the
		// (invalid) authenticatedUserID value from their session and call the next
		// handler in the chain as normal.
		user, err := app.users.Get(app.session.GetInt(r, "authenticatedUserID"))
		if errors.Is(err, models.ErrNoRecord) || !user.Active {
			app.session.Remove(r, "authenticatedUserID")

			next.ServeHTTP(rw, r)

			return
		} else if err != nil {
			app.serverError(rw, err)

			return
		}
		// if we have confirmed that the request is coming from an active, authenticated user,
		// We create a new copy of the request, with a true boolean value
		// added to the request context to indicate this, and call the next handler
		// in the chain *using this new copy of the request*.
		ctx := context.WithValue(r.Context(), contextKeyIsAuthenticated, true)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
