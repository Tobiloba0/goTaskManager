package handlers

import (
	"net/http"
	"strconv"
	"sync"
	"gotask/models"
	"github.com/gin-gonic/gin"
)

var (
	tasks  []models.Task
	nextID uint = 1
	mu     sync.Mutex 
)

func UpdateTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var input models.UpdateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// Find the task and replace its values completely
	for i, t := range tasks {
		if t.ID == uint(id) {
			tasks[i].Title = input.Title
			tasks[i].Description = input.Description
			tasks[i].Status = input.Status

			c.JSON(http.StatusOK, gin.H{"data": tasks[i]})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
}

// Task 2 — Add DeleteTask (DELETE /tasks/:id)
func DeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i, t := range tasks {
		if t.ID == uint(id) {
			// Remove the element from the slice safely
			tasks = append(tasks[:i], tasks[i+1:]...)
			
			// 204 No Content sends no body back
			c.Status(http.StatusNoContent)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
}

// Task 3 — Filter by Status (Updated GetTasks)
func GetTasks(c *gin.Context) {
	statusFilter := c.Query("status") // Reads ?status=pending

	mu.Lock()
	defer mu.Unlock()

	// If no filter query parameter is passed, return everything
	if statusFilter == "" {
		c.JSON(http.StatusOK, gin.H{"data": tasks})
		return
	}

	// Filter tasks dynamically into a temporary slice
	filtered := []models.Task{}
	for _, t := range tasks {
		if string(t.Status) == statusFilter {
			filtered = append(filtered, t)
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": filtered})
}

// Bonus — Add a PatchTask (PATCH /tasks/:id)
func PatchTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var input models.PatchTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i, t := range tasks {
		if t.ID == uint(id) {
			// Check if pointer is not nil; if it holds a address, dereference and apply it
			if input.Title != nil {
				tasks[i].Title = *input.Title
			}
			if input.Description != nil {
				tasks[i].Description = *input.Description
			}
			if input.Status != nil {
				tasks[i].Status = *input.Status
			}

			c.JSON(http.StatusOK, gin.H{"data": tasks[i]})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
}