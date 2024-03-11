package models


type User struct {
	ID 	     int64
	Email    string `binding:"required"`
	Password string `binding:"required"`

}