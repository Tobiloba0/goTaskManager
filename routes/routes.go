package routes
 
import (
    "github.com/gin-gonic/gin"
    "gotask/config"
    "gotask/handlers"
    "gotask/middleware"
)
 
func SetupRoutes(cfg config.Config, r *gin.Engine) {
    api := r.Group("/api/v1")
 
    auth := api.Group("/auth")
    {
        auth.POST("/register", handlers.RegisterUser)
        auth.POST("/login",    handlers.LoginUser)
    }
 
    // Protected routes — RequireAuth middleware runs first
    tasks := api.Group("/tasks")
    tasks.Use(middleware.RequireAuth(cfg.JWTSecret))
    {
        tasks.GET("",      handlers.GetTasks)
        tasks.POST("",     handlers.CreateTask)
        tasks.GET("/:id",  handlers.GetTask)
        tasks.PUT("/:id",  handlers.UpdateTask)
        tasks.DELETE("/:id", handlers.DeleteTask)
    }
}
