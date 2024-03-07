package main

import (
	"github.com/dkr290/go-events-booking-api/db"
	"github.com/dkr290/go-events-booking-api/handlers"
	"github.com/gin-gonic/gin"
)

func main() {

	db.InitDB()

	server := gin.Default()
	handlers.RegisterRoutes(server)

	server.Run(":8080") //localhost:8080 port for listening

}
