package db

import (
	"database/sql"
	"errors"

	"github.com/mkabdelrahman/snippetbox/model"
)

type SnippetStore struct {
	db *sql.DB
}

func NewSnippetStore(db *sql.DB) *SnippetStore {

	return &SnippetStore{
		db: db,
	}
}

func (s *SnippetStore) GetById(id int) (*model.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
WHERE expires > UTC_TIMESTAMP() AND id = ?`
	row := s.db.QueryRow(stmt, id)

	var snippet model.Snippet
	err := row.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &snippet, nil
}

func (s *SnippetStore) Insert(params model.NewSnippetParams) (*model.Snippet, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := s.db.Exec(stmt, params.Title, params.Content, params.Expires)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	snippet, err := s.GetById(int(id))
	if err != nil {
		return nil, err
	}
	return snippet, nil
}
