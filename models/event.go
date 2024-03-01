package models

import (
	"time"

	"github.com/dkr290/go-events-booking-api/db"
)

// our events sturcture and models package will have methods related to deatabase stuff
type Event struct {
	ID          int64
	Name        string    `binding:"required"`
	Description string    `binding:"required"`
	Location    string    `binding:"required"`
	DateTime    time.Time `binding:"required"`
	UserID      int
}

func (e *Event) Save() error {
	// adding to the database

	query := `
	         INSERT INTO events(name,description,location,datetime,user_id)
	         VALUES (?, ?, ?, ?, ?)`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(
		e.Name,
		e.Description,
		e.Location,
		e.DateTime,
		e.UserID)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	e.ID = id
	return err

}

// call all available events
func GetAllEvents() ([]Event, error) {

	query := "SELECT * FROM events"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var events []Event

	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func GetEventById(id int64) (*Event, error) {

	query := "SELECT * FROM EVENTS where id = ?"
	row := db.DB.QueryRow(query, id)
	var event Event
	if err := row.Scan(
		&event.ID,
		&event.Name,
		&event.Description,
		&event.Location,
		&event.DateTime,
		&event.UserID); err != nil {
		return &Event{}, err
	}

	return &event, nil

}

func (e *Event) Update() error {
	query := `
	  UPDATE events
	  SET name = ? , description = ?, location = ? , dateTime = ? 
	  WHERE id = ?
	`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.ID)

	return err
}
