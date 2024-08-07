package models

import (
	"database/sql"
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
	return nil, nil
}

// return top 10 most recent snippets
func (m *SnippetsModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
