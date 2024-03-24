package middlewares

import (
	"net/http"

	
	"github.com/dkr290/events-bookings-new/utils"
	"github.com/gin-gonic/gin"
)

func Authenticate(c *gin.Context){

   token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "not authorized"})
		return
	}

	userId, err := utils.VerifyToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": " unauthorized error",
			"error": err.Error(),

		})
	}

	//set the userId to pass it with the contect if needed in events
	c.Set("userId",userId)

	c.Next()



}