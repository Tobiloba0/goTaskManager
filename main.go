package main
 
import (
    "github.com/gin-gonic/gin"
	"gotask/routes"

)
 
func main() {
    // gin.Default() gives us a router with Logger and Recovery middleware
    r := gin.Default()
 
    // A simple health-check route
   routes.SetupRoutes(r)
 
    // Start the server on port 8080
    r.Run(":8080")
}
