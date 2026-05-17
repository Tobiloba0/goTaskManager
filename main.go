package main

import (
	"github.com/gin-gonic/gin"
	"gotask/routes"
	"net/http"
)

func main() {
	// gin.Default() gives us a router with Logger and Recovery middleware
	r := gin.Default()
	routes.SetupRoutes(r)

	// A simple health-check route
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// Start the server on port 8080
	r.Run(":8080")
}
