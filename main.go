package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gotask/config"
	"gotask/db"
	"gotask/routes"
)

func main() {
	// 1. Load the local .env configuration file
	if err := godotenv.Load(); err != nil {
		log.Println("Notice: No .env file found, reading from system environment variables")
	}

	// 2. Initialize app configuration settings and link database
	cfg := config.Load()
	db.Connect(cfg) // This connects GORM and automatically migrates the tables

	// 3. Fire up the Gin engine router instance
	r := gin.Default()
	
	routes.SetupRoutes(cfg, r)

	// 4. Run the web server using your custom environment port
	log.Printf("Server initializing on port %s...", cfg.Port)
	r.Run(":" + cfg.Port)
}