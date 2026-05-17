package routes

import (
	"github.com/gin-gonic/gin"
	"gotask/handlers"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		tasks := api.Group("/tasks")
		{
			tasks.GET("", handlers.GetTasks)
			tasks.POST("", handlers.CreateTask)
			tasks.GET("/:id", handlers.GetTask)
			tasks.PUT("/:id", handlers.UpdateTask)     // Task 1
			tasks.DELETE("/:id", handlers.DeleteTask)  // Task 2
			tasks.PATCH("/:id", handlers.PatchTask)    // Bonus
		}
	}
}