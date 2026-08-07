package controllers

import (
	"net/http"
	"strconv"

	"task_manager/models"
	"task_manager/services"

	"github.com/gin-gonic/gin"
)

func GetTasks(c *gin.Context) {

	tasks := services.GetAllTasks()

	c.JSON(http.StatusOK, tasks)

}

func GetTask(c *gin.Context) {

	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid task ID",
		})
		return
	}

	task, found := services.GetTaskByID(id)

	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "task not found",
		})
		return
	}

	c.JSON(http.StatusOK, task)
}

func CreateTask(c *gin.Context) {

	var task models.Task

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if !services.IsValidStatus(task.Status) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "status must be Pending, In Progress, or Completed",
		})
		return
	}

	newTask := services.AddTask(task)

	c.JSON(http.StatusCreated, newTask)
}

func UpdateTask(c *gin.Context) {

	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid task ID",
		})
		return
	}

	var updatedTask models.Task

	if err := c.ShouldBindJSON(&updatedTask); err != nil {

		if !services.IsValidStatus(updatedTask.Status) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "status must be Pending, In Progress, or Completed",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	task, found := services.UpdateTask(id, updatedTask)

	if !found {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "task not found",
		})

		return
	}

	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid task ID",
		})
		return

	}
	deleted := services.DeleteTask(id)

	if !deleted {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "task not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "task deleted successfully",
	})

}
