package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/dkr290/go-advanced-projects/go-templ-cruid/models"
	"github.com/go-sql-driver/mysql"
)

type TodoDatabase interface {
	GetAllTasks() ([]models.Task, error)
}

type MysqlDatabase struct {
	DB *sql.DB
}

func InitDB(cfg mysql.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg.FormatDSN())
	count := 0

	retryInterval := 2 * time.Second
	for {
		if err := db.Ping(); err == nil {
			log.Println("Sucesfully connected to the database")
			return db, nil
		} else {
			log.Printf("Attempt %d: Failed to connect to the database. Retrying in %v...\n", count, retryInterval)
			time.Sleep(retryInterval)
			count++
			if count > 10 {
				return nil, err
			}
		}
	}
}

func (d *MysqlDatabase) GetAllTasks() ([]models.Task, error) {
	query := "SELECT id,task,done FROM tasks"

	rows, err := d.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var todo models.Task
		rowErr := rows.Scan(&todo.Id, &todo.Task, &todo.Done)
		if rowErr != nil {
			return nil, rowErr
		}
		tasks = append(tasks, todo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}
