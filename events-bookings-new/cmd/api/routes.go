package main

import (
	"github.com/dkr290/events-bookings-new/db"
	"github.com/dkr290/events-bookings-new/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {

	db := db.MySQLDatabase{}
	h := handlers.New(db)
	h.DB.InitDB()
	h.DB.CreateTables()

	server.GET("/events", h.GetEvents)
	server.POST("/events", h.CreateEvent)
	server.GET("/events/:id", h.GetEvent) ///events/1 , /events/5
	server.PUT("/events/:id", h.UpdateEvent)
	server.DELETE("/events/:id", h.DeleteEvent)
	server.POST("/signup", h.Signup)
	server.POST("/login", h.Login)
}
