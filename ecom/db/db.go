package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dkr290/go-advanced-projects/ecom/types"
	"github.com/go-sql-driver/mysql"
)

type Database interface {
	GetUserByEmail(email string) (*types.User, error)
	GetUserById(id int) (*types.User, error)
	CreateUser(user types.User) error
	InitDB(cfg mysql.Config) (*sql.DB, error)
}

type MysqlDB struct {
	db *sql.DB
}

func (m *MysqlDB) InitDB(cfg mysql.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

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
func (m *MysqlDB) GetUserByEmail(email string) (*types.User, error) {

	rows, err := m.db.Query("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}

	u := new(types.User)
	for rows.Next() {
		err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.CreatedAt)
		if err != nil {
			return nil, err
		}

	}
	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return u, nil
}

func (m *MysqlDB) GetUserById(id int) (*types.User, error) {
	rows, err := m.db.Query("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	u := new(types.User)
	for rows.Next() {
		err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.CreatedAt)
		if err != nil {
			return nil, err
		}

	}
	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return u, nil

}
func (m *MysqlDB) CreateUser(user types.User) error {
	if m.db == nil {
		return errors.New("database connection is nil")
	}
	_, err := m.db.Exec("INSERT INTO users (firstName,lastName,email,password) VALUES(?,?,?,?)",
		user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}
