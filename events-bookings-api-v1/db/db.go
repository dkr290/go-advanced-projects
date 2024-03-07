package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	count := 10
	var err error
	for count > 0 {
		DB, err = sql.Open("sqlite3", "api.db")
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

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	createTables()

}

func createTables() {
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
	_, err := DB.Exec(createEventsTable)
	if err != nil {
		panic("could not create events table")
	}
}
