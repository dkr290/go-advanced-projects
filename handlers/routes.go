package handlers

import "github.com/gin-gonic/gin"

func RegisterRoutes(server *gin.Engine) {
	// the handlers
	server.GET("/events", getEvents)
	server.POST("/events", createEvent)
	server.GET("/events/:id", getEvent) ///events/1 , /events/5
	server.PUT("/events/:id", updateEvent)
}
