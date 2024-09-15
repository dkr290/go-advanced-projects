package userhandlers

import (
	"testing"

	"github.com/dkr290/go-advanced-projects/ecom/types"
	"github.com/go-sql-driver/mysql"
)

func TestUserHandlers(t *testing.T) {
	db := &mockMysqlDB{}
	handler := NewUserHandler(db)
}

type mockMysqlDB struct{}

func (m *mockMysqlDB) GetUserByEmail(email string) (*types.User, error) {
	return nil, nil
}
func (m *mockMysqlDB) GetUserById(id int) (*types.User, error) {
	return nil, nil
}

func (m *mockMysqlDB) CreateUser(user types.User) error {
	return nil
}
func (m *mockMysqlDB) InitDB(cfg mysql.Config) error {
	return nil
}
