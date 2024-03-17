package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/dkr290/events-bookings-new/models"
	"github.com/dkr290/events-bookings-new/utils"
	_ "github.com/mattn/go-sqlite3"
)

type Database interface {
	InitDB()
	CreateTables()
	Save(event models.Event) error
	Update(event models.Event) error
	Delete(event *models.Event) error
	GetAllEvents() ([]models.Event, error)
	GetEventById(id int64) (*models.Event, error)
	SaveUser(u models.User) error
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

	createUsersTable := `
       CREATE TABLE IF NOT EXISTS users(
       	id INTEGER PRIMARY KEY AUTOINCREMENT,
       	email TEXT NOT NULL UNIQUE,
       	password TEXT NOT NULL

       )  
    `
	_, err := m.DB.Exec(createUsersTable)
	if err != nil {
		panic("could not create users table")
	}

	createEventsTable := `
	   CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		datetime DATETIME NOT NULL,
		user_id INTEGER,
		FOREIGN KEY(user_id) REFERENCES users(id)
	)
	`
	_, err = m.DB.Exec(createEventsTable)
	if err != nil {
		panic("could not create events table")
	}

	createResgistrationsTable := `
      CREATE TABLE IF NOT EXISTS registrations (
      	id INTEGER PRIMARY KEY AUTOINCREMENT,
      	event_id INTEGER,
      	user_id INTEGER,
      	FOREIGN KEY(event_id) REFERENCES events(id)
      	FOREIGN KEY(user_id)  REFERENCES users(id)
      )
	`
	_ , err = m.DB.Exec(createResgistrationsTable)
	if err != nil {
        panic("could not create registrati table")
	}
}

func (m *MySQLDatabase) Save(event *models.Event) error {
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

func (m *MySQLDatabase) Update(event models.Event) error {
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

func (m *MySQLDatabase) SaveUser(u models.User) error {
	query := "INSERT INTO users(email,password) VALUES (?, ?)"
	stmt, err := m.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()
	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(
		u.Email,
		hashedPassword)
	if err != nil {
		return err
	}

	userId, err := result.LastInsertId()

	u.ID = userId
	return err //if it dont have error this is nil so exatly what is needed

}

func (m *MySQLDatabase) ValidateCredentials(u *models.User) error {

	query := "SELECT id, password FROM users WHERE email = ?"
	row := m.DB.QueryRow(query, u.Email)

	var retreivedPassword string
	err := row.Scan(&u.ID, &retreivedPassword)

	if err != nil {
		return errors.New("credentials invalid")
	}

	passwordIsValid := utils.CheckPasswordHash(u.Password, retreivedPassword)

	if !passwordIsValid {
		return errors.New("credentials invalid")
	}

	return nil

}

func (m *MySQLDatabase) GetAllUsers() ([]models.User, error) {

	query := "SELECT * FROM Users"
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var allUsers []models.User

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}
		allUsers = append(allUsers, user)
	}

	return allUsers, nil
}
