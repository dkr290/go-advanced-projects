package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/dkr290/go-advanced-projects/pic-dream-api/pkg/db"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func createDB() (*sql.DB, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	host := os.Getenv("DB_HOST")
	if len(host) == 0 {
		log.Fatal("DB_HOST is mandatory")
	}
	user := os.Getenv("DB_USER")
	if len(user) == 0 {
		log.Fatal("DB_USER is mandatory")
	}
	pass := os.Getenv("DB_PASSWORD")
	if len(pass) == 0 {
		log.Fatal("DB_PASSWORD is mandatory")
	}
	dbname := os.Getenv("DB_NAME")
	if len(dbname) == 0 {
		log.Fatal("DB_NAME is mandatory")
	}

	return db.CreateDatabase(dbname, user, pass, host)
}

func main() {
	db, err := createDB()
	if err != nil {
		log.Fatal(err)
	}

	tables := []string{
		"schema_migrations",
		"accounts",
		"images",
	}

	for _, table := range tables {
		query := fmt.Sprintf("drop table if exists %s cascade", table)
		if _, err := db.Exec(query); err != nil {
			log.Fatal(err)
		}
	}
}
