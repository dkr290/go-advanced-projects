package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func(h *Handlers) RegisterForEvent(c *gin.Context){
	userId := c.GetInt64("userId")
	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not get event id " + err.Error()})
		return
	}

   
   event, err := h.DB.GetEventById(eventId)
   if err != nil {
   	c.JSON(http.StatusInternalServerError,gin.H{
   		"message": "could not fetch event",
   	})
   }

   err = h.DB.Register(userId,event)
   if err != nil {
   	c.JSON(http.StatusInternalServerError,gin.H{
   		"message": "could not register the event",
   	})
   }

   c.JSON(http.StatusCreated,gin.H{"message":"registered"})
}

func(h *Handlers) CancelRegistration(c *gin.Context){}