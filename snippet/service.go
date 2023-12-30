package snippet

import (
	"database/sql"
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

type Service struct {
	DB *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{
		DB: db,
	}
}

// Insert will insert a new snippet into the database.
func (m *Service) Insert(params NewModelParams, userTimeZone string) (int, error) {
	// Combine date and time into a single time.Time object.
	expires := time.Date(
		params.ExpiresDate.Year(),
		params.ExpiresDate.Month(),
		params.ExpiresDate.Day(),
		params.ExpiresTime.Hour(),
		params.ExpiresTime.Minute(),
		params.ExpiresTime.Second(),
		0, // Nanoseconds
		params.ExpiresTime.Location(),
	)

	// Convert expires time to UTC before inserting into the database.
	expiresUTC := expires.UTC()

	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), ?)`

	result, err := m.DB.Exec(stmt, params.Title, params.Content, expiresUTC)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Get will return a specific snippet based on its id.
func (m *Service) Get(id int, userTimeZone string) (*Model, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt, id)
	s := &Model{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	// Convert expires time to user's time zone before returning.
	s.Expires = s.Expires.In(time.FixedZone(userTimeZone, 0))

	return s, nil
}

// Delete will delete a specific snippet based on its id.
func (m *Service) Delete(id int) error {
	stmt := `DELETE FROM snippets WHERE id = ?`

	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}

// Latest will return the 10 most recently created snippets.
func (m *Service) Latest(userTimeZone string) ([]*Model, error) {
	// Write the SQL statement we want to execute.
	stmt := `SELECT id, title, content, created, expires FROM snippets
		WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	snippets := []*Model{}

	for rows.Next() {
		s := &Model{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		// Convert expires time to user's time zone before appending to the result.
		s.Expires = s.Expires.In(time.FixedZone(userTimeZone, 0))

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}

// Delete prepares the DELETE statement for a specific snippet.
func (s *Model) Delete() (string, []interface{}) {
	stmt := `DELETE FROM snippets WHERE id = ?`
	params := []interface{}{s.ID}

	return stmt, params
}
