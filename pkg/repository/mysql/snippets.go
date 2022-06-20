package mysql

import (
	"database/sql"
	"errors"

	"github.com/mycok/snippet-bin/pkg/models"
)

// Perform a compile time / static check to verify whether *SnippetModel type
// satisfies UserRepository interface.
var _ models.SnippetRepository = (*SnippetModel)(nil)

const (
	insertSnippet = `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	getSnippetByID = `SELECT id, title, content, created, expires
					FROM snippets
					WHERE expires > UTC_TIMESTAMP() and id = ?`

	latestSnippets = `SELECT id, title, content, created, expires
					FROM snippets
					WHERE expires > UTC_TIMESTAMP()
					ORDER BY created
					DESC
					LIMIT 10`
)

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	result, err := m.DB.Exec(insertSnippet, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	s := &models.Snippet{}

	err := m.DB.QueryRow(getSnippetByID, id).Scan(
		&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	rows, err := m.DB.Query(latestSnippets)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []*models.Snippet{}

	for rows.Next() {
		s := &models.Snippet{}

		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
