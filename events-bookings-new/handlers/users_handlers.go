package handlers

import (
	"net/http"

	"github.com/dkr290/events-bookings-new/models"
	"github.com/dkr290/events-bookings-new/utils"
	"github.com/gin-gonic/gin"
)

func (h *Handlers) Signup(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not parse request data " + err.Error()})
		return

	}

	err := h.DB.SaveUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not save user",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created sucessfully",
		"email":   user.Email,
	})

}

func (h *Handlers) Login(c *gin.Context) {

	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "could not parse request data ",
			"error": err.Error(),
		})
		return

	}

	err := h.DB.ValidateCredentials(user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "clould not authenticate user",
			"error":   err.Error(),
		})
		return
	}

	token, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not authenticate user",
			"error":   err.Error(),
		})

	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login sucessfull",
		"token":   token,
	})
}

func (h *Handlers) GetUsers(c *gin.Context) {

	users, err := h.DB.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not fetch users. Try again later.",
		})
		return
	}
	c.JSON(http.StatusOK, users)

}
