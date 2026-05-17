package handlers
 
import (
    "net/http"
    "strconv"
 
    "github.com/gin-gonic/gin"
    "gotask/db"
    "gotask/models"
)
 
func GetTasks(c *gin.Context) {
    var tasks []models.Task

    // Fetch the logged-in user from the JWT middleware context
    user, _ := c.MustGet("currentUser").(models.User)

    // 1. Establish defaults if query params aren't explicitly declared
    pageStr := c.DefaultQuery("page", "1")
    limitStr := c.DefaultQuery("limit", "10")

    // 2. Convert strings safely to integers
    page, err := strconv.Atoi(pageStr)
    if err != nil || page < 1 {
        page = 1
    }

    limit, err := strconv.Atoi(limitStr)
    if err != nil || limit < 1 {
        limit = 10
    }

    // 3. Compute structural pagination offset
    offset := (page - 1) * limit

    // 4. Build query chains dynamically, scoped specifically to this user!
    query := db.DB.Model(&models.Task{}).Where("user_id = ?", user.ID)
    
    if status := c.Query("status"); status != "" {
        query = query.Where("status = ?", status)
    }

    // 5. Query the database using GORM Limit and Offset directives
    if err := query.Limit(limit).Offset(offset).Find(&tasks).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "page":  page,
        "limit": limit,
        "data":  tasks,
    })
}
 
func CreateTask(c *gin.Context) {
    var input models.CreateTaskInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    user, _ := c.MustGet("currentUser").(models.User)
    task := models.Task{
        UserID:      user.ID,
        Title:       input.Title, 
        Description: input.Description,
        DueDate:     input.DueDate, 
    }
    
    if err := db.DB.Create(&task).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, gin.H{"data": task})
}
 
func GetTask(c *gin.Context) {
    var task models.Task
    user, _ := c.MustGet("currentUser").(models.User)

    // Find the task by ID AND make sure it belongs to the logged-in user
    if err := db.DB.Where("id = ? AND user_id = ?", c.Param("id"), user.ID).First(&task).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "task not found or access denied"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": task})
}
 
func UpdateTask(c *gin.Context) {
    var task models.Task
    user, _ := c.MustGet("currentUser").(models.User)

    // Secure the lookup so you can only fetch your own task
    if err := db.DB.Where("id = ? AND user_id = ?", c.Param("id"), user.ID).First(&task).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "task not found or access denied"})
        return
    }
 
    var input models.UpdateTaskInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
 
    db.DB.Model(&task).Updates(input)
    c.JSON(http.StatusOK, gin.H{"data": task})
}
 
func DeleteTask(c *gin.Context) {
    var task models.Task
    user, _ := c.MustGet("currentUser").(models.User)

    // Secure the lookup so you can only delete your own task
    if err := db.DB.Where("id = ? AND user_id = ?", c.Param("id"), user.ID).First(&task).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "task not found or access denied"})
        return
    }

    db.DB.Delete(&task)
    c.JSON(http.StatusNoContent, nil) // 204 No Content success
}