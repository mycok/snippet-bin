package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golangcollege/sessions"
)

func Test_homeHandler(t *testing.T) {
	testCases := []struct {
		name               string
		templatePath string
		expectedStatusCode int
	}{
		{
			name:               "WithDBError",
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:               "WithMissingTemplate",
			templatePath: "./ui/html/",
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:               "WithSnippets",
			templatePath: "../../ui/html/",
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app := &application{
				snippets: &handlerSnippetRepositoryMock{
					testName: tc.name,
				},
				errLog: log.New(io.Discard, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
			}

			rr := httptest.NewRecorder()

			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			ctx := context.WithValue(req.Context(), contextKeyIsAuthenticated, true)
			req = req.WithContext(ctx)

			tempCache, err := newTemplateCache(tc.templatePath)
			if err != nil {
				t.Fatal(err)
			}

			app.templateCache = tempCache

			// Create a new session object to add to the request context.
			session := sessions.New([]byte("yeueuu+hffs24453+42fggsg*yu@etyr"))
			session.Lifetime = 10 * time.Second

			app.session = session

			app.session.Enable(http.HandlerFunc(app.home)).ServeHTTP(rr, req)

			resp := rr.Result()

			if tc.expectedStatusCode != resp.StatusCode {
				t.Errorf(
					"Expected status code: %d, but got: %d instead",
					tc.expectedStatusCode,
					resp.StatusCode,
				)
			}
		})
	}
}
