package main

import (
	"github.com/dkr290/events-bookings-new/db"
	"github.com/dkr290/events-bookings-new/handlers"
	"github.com/dkr290/events-bookings-new/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {

	db := db.MySQLDatabase{}
	h := handlers.New(db)
	h.DB.InitDB()
	h.DB.CreateTables()

	server.GET("/events", h.GetEvents)
	server.GET("/events/:id", h.GetEvent) ///events/1 , /events/5

	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate)
	authenticated.POST("/events", h.CreateEvent)
	authenticated.PUT("/events/:id", h.UpdateEvent)
	authenticated.DELETE("/events/:id", h.DeleteEvent)


	server.POST("/signup", h.Signup)
	server.POST("/login", h.Login)

//	server.GET("/users", h.GetUsers)

}
