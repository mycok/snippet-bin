package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate mail")
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Active         bool
}

type SnippetRepository interface {
	Insert(title, content, expires string) (int, error)
	Get(id int) (*Snippet, error)
	Latest() ([]*Snippet, error)
}

type UserRepository interface {
	Insert(name, email, password string) error
	Get(id int) (*User, error)
	Authenticate(email, password string) (int, error)
}
