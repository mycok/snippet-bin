package models

import (
	"errors"
	"time"
)

var ErrNorRecord = errors.New("models: no matching record found")

type Snippet struct {
	ID int
	Tittle string
	Content string
	Created time.Time
	Expires time.Time
}
