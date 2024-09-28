package db

import (
	"database/sql"
	"fmt"

	"github.com/dkr290/go-advanced-projects/ecom/types"
)

type UserDatabaseInt interface {
	GetUserByEmail(email string) (*types.User, error)
	GetUserById(id int) (*types.User, error)
	CreateUser(user types.User) error
}

type UserMysqlDB struct {
	DB *sql.DB
}

func (m *UserMysqlDB) GetUserByEmail(email string) (*types.User, error) {

	rows, err := m.DB.Query("SELECT * FROM users WHERE email = ?", email)
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

func (m *UserMysqlDB) GetUserById(id int) (*types.User, error) {
	rows, err := m.DB.Query("SELECT * FROM users WHERE id = ?", id)
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
func (m *UserMysqlDB) CreateUser(user types.User) error {
	_, err := m.DB.Exec("INSERT INTO users (firstName,lastName,email,password) VALUES(?,?,?,?)",
		user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}
