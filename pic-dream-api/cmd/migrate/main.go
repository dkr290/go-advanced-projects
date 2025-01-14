package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/dkr290/go-advanced-projects/pic-dream-api/pkg/db"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
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

	// Create migration instance
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Point to your migration files. Here we're using local files, but it could be other sources.
	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations", // source URL
		"postgres",                      // database name
		driver,                          // database instance
	)
	if err != nil {
		log.Fatal(err)
	}

	cmd := os.Args[len(os.Args)-1]
	if cmd == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
	if cmd == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
}
