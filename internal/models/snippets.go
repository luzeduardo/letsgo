package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

//Define a SnippetModel type thar wrapps a sql.DB connection pool

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	//ensures sql.Rows results is always properly closed before the Latest() method returns
	// must become after the error check from Query() method
	// if Query returns an error, it'll result in a panic trying to close a
	// nil resultset
	defer rows.Close()

	snippets := []*Snippet{}

	//iterate through the rows in the resultset
	//it closes automatically and closes itself and frees-up the DB connection
	for rows.Next() {
		s := &Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	//after rows.Next ends, checks if some error ocurred during the iteration
	// do not assume that the iteration completed successfully
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires 
	FROM
	 snippets WHERE expires > UTC_TIMESTAMP() AND id = ?`
	row := m.DB.QueryRow(stmt, id)
	// Initialize a pointer to the new struct
	s := &Snippet{}

	// copy values from each field in sql.Row to the corresponding field in the Snippet struct
	// the args are *pointers* to the place I want to copy the data into
	// the number of args must be exactly the same as the number of columns returned by the statement
	// convert the raw output of the SQL to the native Go types
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// if the query return no rows, the row.Scan() will return a sql.ErrNoROws error.
		// use the errors.Is to check for a specifically error and return own ErrNoRecord error
		// instead
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	//return the Snippet struct
	return s, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
