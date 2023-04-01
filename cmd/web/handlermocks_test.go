package main

import (
	"bytes"
	"fmt"
	"html/template"
	"testing"

	"github.com/mycok/snippet-bin/pkg/models"
)

// Mock the SnippetRepository.
type handlerSnippetRepositoryMock struct {
	testName string
}

func (r *handlerSnippetRepositoryMock) Insert(title, content, expires string) (int, error) {
	return 0, nil
}

func (r *handlerSnippetRepositoryMock) Get(id int) (*models.Snippet, error) {
	return &models.Snippet{}, nil
}

func (r *handlerSnippetRepositoryMock) Latest() ([]*models.Snippet, error) {
	if r.testName == "WithDBError" {
		return nil, fmt.Errorf("Database error")
	}

	return []*models.Snippet{}, nil
}

func getParsedTemplate(t *testing.T, name string, td templateData) []byte {
	t.Helper()

	temp, err := getCachedTemplate(name)
	if err != nil {
		t.Fatal(err)
	}

	tempBuf := new(bytes.Buffer)

	// Parse the template with the provided data.
	if err := temp.Execute(tempBuf, td); err != nil {
		t.Fatal(err)
	}

	return tempBuf.Bytes()
}

func getCachedTemplate(name string) (*template.Template, error) {
	tempCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		return nil, err
	}

	temp, ok := tempCache[name]
	if !ok {
		return nil, fmt.Errorf("the template %s does not exit", name)
	}

	return temp, nil
}
