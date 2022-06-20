package main

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mycok/snippet-bin/pkg/models"
)

// Mock the UserRepository.
type mockUserRepository struct {
	testName string
}

func (r *mockUserRepository) Insert(name, email, password string) error {
	return nil
}

func (r *mockUserRepository) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (r *mockUserRepository) Get(id int) (*models.User, error) {
	switch r.testName {
	case "WithNonExistingAuthID":
		return nil, models.ErrNoRecord
	case "WithAuthIDForDeactivatedUser":
		return &models.User{Active: false}, nil
	case "WithAuthIDAndDBError":
		return &models.User{Active: true}, errors.New("this should be a database error")
	case "WithValidAuthID":
		return &models.User{Active: true}, nil
	default:
		return nil, nil
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
	okHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(http.StatusText(http.StatusOK)))
	})

	// Create a request object to pass to the middleware.
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	return rr, req, okHandler
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
