package main

import (
	"net/http"

	"github.com/dkr290/go-events-booking-api/db"
	"github.com/dkr290/go-events-booking-api/models"
	"github.com/gin-gonic/gin"
)

func main() {

	db.InitDB()

	server := gin.Default()

	//the handlers
	server.GET("/events", getEvents)
	server.POST("/events", createEvent)

	server.Run(":8080") //localhost:8080 port for listening

}

func getEvents(c *gin.Context) {

	events, err := models.GetAllEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not fetch events. Try again later.",
		})
		return
	}
	c.JSON(http.StatusOK, events)

}

func createEvent(c *gin.Context) {
	var event = models.Event{}
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not parse request data" + err.Error()})
		return

	}

	//use some dummy value
	event.ID = 1
	event.UserID = 1

	err := event.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not save events. Try again later.",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Event created",
		"event":   event,
	})

}
