package handlers

import (
	"net/http"
	"strconv"

	"github.com/dkr290/events-bookings-new/db"
	"github.com/dkr290/events-bookings-new/models"
	"github.com/dkr290/events-bookings-new/utils"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	DB db.MySQLDatabase
}

func New(db db.MySQLDatabase) *Handlers {
	return &Handlers{
		DB: db,
	}
}

func (h *Handlers) GetEvents(c *gin.Context) {

	events, err := h.DB.GetAllEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not fetch events. Try again later.",
		})
		return
	}
	c.JSON(http.StatusOK, events)

}

func (h *Handlers) CreateEvent(c *gin.Context) {

	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "not authorized"})
		return
	}

	userId, err := utils.VerifyToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": " unauthorized"})
	}

	var event = models.Event{}
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not parse request data " + err.Error()})
		return

	}

	//this should be the used id of the user who did the event
	event.UserID = userId

	err = h.DB.Save(&event)
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

func (h *Handlers) GetEvent(c *gin.Context) {

	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not get event id " + err.Error()})
		return
	}

	event, err := h.DB.GetEventById(eventId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not fetch event " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, event)

}

func (h *Handlers) UpdateEvent(c *gin.Context) {
	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not get event id " + err.Error()})
		return
	}

	_, err = h.DB.GetEventById(eventId)
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

	if err := h.DB.Update(updatedEvent); err != nil {
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

func (h *Handlers) DeleteEvent(c *gin.Context) {
	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not get event id " + err.Error()})
		return
	}
	event, err := h.DB.GetEventById(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch event " + err.Error()})
		return
	}

	if err := h.DB.Delete(event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete event" + err.Error()})
		return

	}

	c.JSON(http.StatusOK, gin.H{"message": "event deleted sucessfully"})

}
