package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_secureHeaders(t *testing.T) {
	rr, req, nextHandler := httpMiddlewareTestSetup(t)

	// Execute the middleware with the mock next handler, request and response
	// recorder.
	secureHeaders(nextHandler).ServeHTTP(rr, req)

	// Call the Result() method on the http.ResponseRecorder to read the results
	// of the test.
	resp := rr.Result()

	xssProtection := resp.Header.Get("x-XSS-Protection")
	if xssProtection != "1; mode=block" {
		t.Errorf(
			"Expected xssProtection header: '1; mode=block', but got: %q instead",
			xssProtection,
		)
	}

	xFrameOptions := resp.Header.Get("X-Frame-Options")
	if xFrameOptions != "deny" {
		t.Errorf(
			"Expected xFrameOptions header: 'deny', but got: %q instead",
			xFrameOptions,
		)
	}

	// Run common assertions on the response object from *http.ResponseRecorder.
	assertOnResponse(t, rr)
}

func Test_logRequest(t *testing.T) {
	rr, req, nextHandler := httpMiddlewareTestSetup(t)

	logBuf := bytes.Buffer{}

	// Create an instance of log.Logger to pass to the application instance.
	infoLog := log.New(&logBuf, "INFO\t", log.Ldate|log.Ltime)

	// Create an instance of application with infoLog as the logger.
	app := &application{
		infoLog: infoLog,
	}

	// Execute the middleware with the mock next handler, request and response
	// recorder.
	app.logRequest(nextHandler).ServeHTTP(rr, req)

	expectedLogStr := fmt.Sprintf(
		"%s - %s %s %s",
		req.RemoteAddr, req.Proto, req.Method, req.URL.RequestURI(),
	)

	if !strings.Contains(logBuf.String(), expectedLogStr) {
		t.Errorf(
			"Expected request log: %s, but got: %s instead",
			expectedLogStr,
			logBuf.String(),
		)
	}

	// Run common assertions on the response object from *http.ResponseRecorder.
	assertOnResponse(t, rr)
}

func Test_recoverFromPanic(t *testing.T) {
	rr, req, _ := httpMiddlewareTestSetup(t)

	expectedErr := "this handler is mean't to panic"

	// Mock a panic handler to be called next in the middleware chain
	//  after the middleware.
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic(expectedErr)
	})

	logBuf := bytes.Buffer{}

	errLog := log.New(&logBuf, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{ errLog: errLog }

	// Execute the middleware with the mock next handler, request and response
	// recorder.
	app.recoverFromPanic(nextHandler).ServeHTTP(rr, req)

	// Call the Result() method on the http.ResponseRecorder to read the results
	// of the test.
	resp := rr.Result()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf(
			"Expected status code: %q, but got: %q instead",
			http.StatusInternalServerError,
			resp.StatusCode,
		)
	}

	if resp.Header.Get("Connection") != "Close" {
		t.Errorf(
			"Expected response header 'Connection': %s, but got: %s instead",
			"Close",
			resp.Header.Get("Connection"),
		)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %q", err)
	}

	defer resp.Body.Close()

	if string(bytes.TrimSpace(body)) != http.StatusText(http.StatusInternalServerError) {
		t.Errorf(
			"Expected error: %s, but got: %s instead.",
			http.StatusText(http.StatusInternalServerError),
			string(body),
		)
	}
}

func Test_requireAuthentication(t *testing.T) {
	testCases := []struct{
		name string
		withAuth bool
		headers map[string][]string
	} {
		{
			name: "NoAuth",
			withAuth: false,
			headers: map[string][]string{
				"Content-Type": {"text/html; charset=utf-8"},
			},
		},
		{
			name: "WithAuth",
			withAuth: true,
			headers: map[string][]string{
				"Content-Type": {"text/plain; charset=utf-8"},
				"Cache-Control": {"no-store"},
			},
		},
	}

	for _, tc := range testCases {
		rr, req, nextHandler := httpMiddlewareTestSetup(t)

		app := &application{}

		t.Run(tc.name, func(t *testing.T) {
			if tc.withAuth {
				// Add authentication to the request.
				ctx := context.WithValue(req.Context(), contextKeyIsAuthenticated, true)
				// Execute the middleware with the mock next handler, request and response
				// recorder.
				app.requireAuthentication(nextHandler).ServeHTTP(rr, req.WithContext(ctx))

				resp := rr.Result()

				expectedStatusCode := http.StatusOK
				
				if expectedStatusCode != resp.StatusCode {
					t.Errorf(
						"Expected status code: %d, but got: %d instead",
						expectedStatusCode,
						resp.StatusCode,
					)
				}

				if tc.headers["Content-Type"][0] != resp.Header.Get("Content-Type") {
					t.Errorf(
						"Expected Content-Type: %s, but got: %s instead",
						tc.headers["Content-Type"][0],
						resp.Header.Get("Content-Type"),
					)
				}

				if tc.headers["Cache-Control"][0] != resp.Header.Get("Cache-Control") {
					t.Errorf(
						"Expected Content-Type: %s, but got: %s instead",
						tc.headers["Cache-Control"][0],
						resp.Header.Get("Cache-Control"),
					)
				}
			} else {
				// Execute the middleware with the mock next handler, request and response
				// recorder.
				app.requireAuthentication(nextHandler).ServeHTTP(rr, req)

				resp := rr.Result()

				expectedStatusCode := http.StatusSeeOther
				
				if expectedStatusCode != resp.StatusCode {
					t.Errorf(
						"Expected an HTTP redirect code: %d, but got: %d instead",
						expectedStatusCode,
						resp.StatusCode,
					)
				}
			}
		})
	}
}

// Setup and return objects of type [*httptest.ResponseRecorder, *http.Request
// http.Handler]. The http.Handler writes []byte("OK") to http.responseWriter.
func httpMiddlewareTestSetup(t *testing.T) (
	*httptest.ResponseRecorder, *http.Request, http.Handler,
) {
	t.Helper()

	rr := httptest.NewRecorder()

	// Mock a handler to be called next in the middleware chain
	//  after the middleware.
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Create a request object to pass to the middleware.
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	return rr, req, nextHandler
}

// Make assertions on the response such as checking
//  [response.statusCode and response.body].
func assertOnResponse(t *testing.T, rr *httptest.ResponseRecorder) {
	// Call the Result() method on the http.ResponseRecorder to read the results
	// of the test.
	resp := rr.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf(
			"Expected status code: %q, but got: %q instead",
			http.StatusOK,
			resp.StatusCode,
		)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %q", err)
	}

	defer resp.Body.Close()

	if string(body) != "OK" {
		t.Errorf("Expected body to have value 'OK', but got: %q", string(body))
	}
}
