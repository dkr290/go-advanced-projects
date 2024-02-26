package db

import (
	"database/sql"
	"fmt"

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
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

}
