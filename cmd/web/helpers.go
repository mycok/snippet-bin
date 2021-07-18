package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/justinas/nosurf"
)

func (app *application) serverError(rw http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errLog.Output(2, trace)

	http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(rw http.ResponseWriter, status int) {
	http.Error(rw, http.StatusText(status), status)
}

// wrapper around clientError for notFoundError
func (app *application) notFoundError(rw http.ResponseWriter) {
	app.clientError(rw, http.StatusNotFound)
}

func (app *application) addDefaultTemplateData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CurrentYear = time.Now().Year()
	td.Flash = app.session.PopString(r, "flash")
	td.IsAuthenticated = app.isAuthenticated(r)
	td.CSRFToken = nosurf.Token(r)

	return td
}

func (app *application) render(rw http.ResponseWriter, r *http.Request, name string, td *templateData) {
	templateSet, ok := app.templateCache[name]
	if !ok {
		app.serverError(rw, fmt.Errorf("The template %s does not exit", name))

		return
	}
	// initialize a new buffer
	// Write the template to the buffer, instead of straight to the
	// http.ResponseWriter. If there's an error, call our serverError helper and then
	// return.
	buf := new(bytes.Buffer)
	// excute the template set passing in dynamic data with the current year injected
	err := templateSet.Execute(buf, app.addDefaultTemplateData(td, r))
	if err != nil {
		app.serverError(rw, err)

		return
	}
	// write the contents of the buffer to http.ResponseWriter
	buf.WriteTo(rw)
}

func (app *application) isAuthenticated(r *http.Request) bool {
	return app.session.Exists(r, "authenticatedUserID")
}