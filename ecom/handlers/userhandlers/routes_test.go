package userhandlers

import "testing"

func TestUserHandlers(t *testing.T){
	db := &mockMysqlDB{} 
	handler := NewUserHandler()
}


type mockMysqlDB struct{}

func(m *mockMysqlDB)
