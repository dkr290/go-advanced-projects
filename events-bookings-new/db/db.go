package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/dkr290/events-bookings-new/models"
	_ "github.com/mattn/go-sqlite3"
)

type Database interface {
	InitDB()
	CreateTables()
	Save() error
	Update() error
	Delete() error
	GetAllEvents() ([]models.Event, error)
	GetEventById(id int64) (*models.Event, error)
}

type MySQLDatabase struct {
	DB *sql.DB
	
}

func (m *MySQLDatabase) InitDB() {
	count := 10
	var err error
	for count > 0 {
		m.DB, err = sql.Open("sqlite3", "api.db")
		if err != nil && count > 0 {
			fmt.Println("cannot connect to the database")
			count -= 1
		}
		if err != nil && count <= 0 {
			panic("cannot connect to the database")
		}
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	m.DB.SetMaxOpenConns(10)
	m.DB.SetMaxIdleConns(5)

	m.CreateTables()

}

func (m MySQLDatabase) CreateTables() {
	createEventsTable := `
	   CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		datetime DATETIME NOT NULL,
		user_id INTEGER
	   )
	`
	_, err := m.DB.Exec(createEventsTable)
	if err != nil {
		panic("could not create events table")
	}
}

func (m *MySQLDatabase) Save(event models.Event) error {
	// adding to the database

	query := `
	         INSERT INTO events(name,description,location,datetime,user_id)
	         VALUES (?, ?, ?, ?, ?)`

	stmt, err := m.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(
		event.Name,
		event.Description,
		event.Location,
		event.DateTime,
		event.UserID)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	event.ID = id
	return err

}

// call all available events
func (m *MySQLDatabase) GetAllEvents() ([]models.Event, error) {

	query := "SELECT * FROM events"
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var events []models.Event

	for rows.Next() {
		var event models.Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func (m *MySQLDatabase) GetEventById(id int64) (*models.Event, error) {

	query := "SELECT * FROM EVENTS where id = ?"
	row := m.DB.QueryRow(query, id)
	var event models.Event
	if err := row.Scan(
		&event.ID,
		&event.Name,
		&event.Description,
		&event.Location,
		&event.DateTime,
		&event.UserID); err != nil {
		return &models.Event{}, err
	}

	return &event, nil

}

func(m *MySQLDatabase) Update(event models.Event) error {
	query := `
	  UPDATE events
	  SET name = ? , description = ?, location = ? , dateTime = ? 
	  WHERE id = ?
	`

	stmt, err := m.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(event.Name, event.Description, event.Location, event.DateTime, event.ID)

	return err
}

func (m *MySQLDatabase) Delete(event *models.Event) error {

	query := "DELETE FROM events WHERE id =?"

	stmt, err := m.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(event.ID)
	return err

}
