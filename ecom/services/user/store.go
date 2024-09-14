package user

import (
	"database/sql"
	"fmt"

	"github.com/dkr290/go-advanced-projects/ecom/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return Store{
		db: db,
	}
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {

	rows, err := s.db.Query("SELECT * FROM users WHERE email = ?", email)
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

func (s *Store) GetUserById(id int) (*types.User, error) {
	return nil, nil
}
func (s *Store) CreateUser(types.User) error {
	return nil
}
