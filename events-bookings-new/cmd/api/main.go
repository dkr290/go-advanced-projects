package main

import (
	
	"github.com/gin-gonic/gin"
)



func main(){

server := gin.Default()
RegisterRoutes(server)


server.Run(":8080") //localhost:8080 port for listening



	
}