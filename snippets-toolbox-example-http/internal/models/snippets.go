package models

import (
	"database/sql"
	"errors"
	"time"
)

//define snipper type and db model that will hold the data for
//individual model
//the fields correspond to mysql snippets table exactly

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expired time.Time
}

// define snippetsMode type that will wrap a sql.DB connection pool

type SnippetsModel struct {
	DB *sql.DB
}

// insert new snippet to the database
func (m *SnippetsModel) Insert(title string, content string, expires int) (int, error) {
	// Write the SQL statement we want to execute
	// for readability (which is why it's surrounded with backquotes instead
	// of normal double quotes).
	stmt := `INSERT INTO snippets (title,content,created,expires) 
	       VALUES (?, ?,UTC_TIMESTAMP(),DATE_ADD(UTC_TIMESTAMP(),INTERVAL ? DAY))`

	//Exec() to execute the statement to sql

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}
	//LastInsertedId to get the ID of the newly inserted record

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err

	}
	//return the id in int

	return int(id), nil
}

// This will return a specific snippet based ont its id
func (m *SnippetsModel) Get(id int) (*Snippet, error) {
	//statement to execute in this case

	stmt := `select id,title,content,created,expires FROM snippets 
	         WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt, id)

	//Initialize a pointer to a new zeroed Snippet struct.

	s := &Snippet{}

	//with row.Scan() copy values from each field into sql.Row to the corresponding fieldin the snippets struct

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expired)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

// return top 10 most recent snippets
func (m *SnippetsModel) Latest() ([]*Snippet, error) {
	stmt := `SELECT id,title,content,created,expires FROM snippets 
	          WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []*Snippet{}

	for rows.Next() {
		s := &Snippet{}
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expired)
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
