package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/dkr290/go-advanced-projects/go-templ-cruid-microservice/backend/models"
	"github.com/go-sql-driver/mysql"
)

type TodoDatabase interface {
	GetAllTasks() ([]models.Task, error)
	AddTask(string) error
	GetTaskByID(int) (*models.Task, error)
	UpdateTaskByID(models.Task) error
	DeleteTaskByID(int) error
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

func (d *MysqlDatabase) AddTask(task string) error {
	query := "INSERT INTO tasks (task) VALUES (?)"
	stmt, err := d.DB.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(task)
	if err != nil {
		return err
	}
	stmt.Close()
	return nil
}

func (d *MysqlDatabase) GetTaskByID(id int) (*models.Task, error) {
	query := "SELECT id, task, done FROM tasks WHERE id = ?"
	var task models.Task
	row := d.DB.QueryRow(query, id)
	err := row.Scan(&task.Id, &task.Task, &task.Done)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("No task was found with id %d", id)
		}
		return nil, err
	}

	return &task, nil
}

func (d *MysqlDatabase) UpdateTaskByID(task models.Task) error {
	query := "UPDATE tasks SET task = ?, done = ? WHERE id = ?"
	result, err := d.DB.Exec(query, task.Task, task.Done, task.Id)
	if err != nil {
		return err
	}
	rowsAffecter, err := result.RowsAffected()
	if rowsAffecter == 0 {
		return fmt.Errorf("no rows affected error:  %v", err)
	}
	if err != nil {
		return err
	}
	return nil
}

func (d *MysqlDatabase) DeleteTaskByID(id int) error {
	query := "DELETE FROM tasks WHERE id = ?"
	stmt, err := d.DB.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	defer stmt.Close()
	return nil
}
