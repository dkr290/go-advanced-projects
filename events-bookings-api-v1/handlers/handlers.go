package handlers

import (
	"net/http"
	"strconv"

	"github.com/dkr290/go-events-booking-api/models"
	"github.com/gin-gonic/gin"
)

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not parse request data " + err.Error()})
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

func getEvent(c *gin.Context) {

	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not get event id " + err.Error()})
		return
	}

	event, err := models.GetEventById(eventId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not fetch event " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, event)

}

func updateEvent(c *gin.Context) {
	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not get event id " + err.Error()})
		return
	}

	_, err = models.GetEventById(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch event " + err.Error()})
		return
	}

	var updatedEvent models.Event
	if err := c.ShouldBindJSON(&updatedEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not parse request data " + err.Error()})
		return

	}

	updatedEvent.ID = eventId

	if err := updatedEvent.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not update event. Try again later.",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Event updated sucessfully",
	})
}

func deleteEvent(c *gin.Context) {
	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not get event id " + err.Error()})
		return
	}
	event, err := models.GetEventById(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch event " + err.Error()})
		return
	}

	if err := event.Delete(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete event" + err.Error()})
		return

	}

	c.JSON(http.StatusOK, gin.H{"message": "event deleted sucessfully"})

}
