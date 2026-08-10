package controllers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"task_manager/data"
	"task_manager/middleware"
	"task_manager/models"
	"task_manager/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type TaskController struct {
	Service     *data.TaskService
	UserService *data.UserService
}

func (tc *TaskController) GetTasks(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	tasks, err := tc.Service.GetAllTasks(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (tc *TaskController) GetTask(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	task, err := tc.Service.GetTaskByID(ctx, id)
	if errors.Is(err, mongo.ErrNoDocuments) {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (tc *TaskController) CreateTask(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var task models.Task

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !services.IsValidStatus(task.Status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status must be Pending, In Progress, or Completed"})
		return
	}

	createdTask, err := tc.Service.CreateTask(ctx, task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create a task"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "task created successfully",
		"task":    createdTask,
	})
}

func (tc *TaskController) UpdateTask(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	var updatedTask models.Task

	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !services.IsValidStatus(updatedTask.Status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status must be Pending, In Progress, or Completed"})
		return
	}

	task, err := tc.Service.UpdateTask(ctx, id, updatedTask)
	if errors.Is(err, mongo.ErrNoDocuments) {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (tc *TaskController) DeleteTask(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	err = tc.Service.DeleteTask(ctx, id)
	if errors.Is(err, mongo.ErrNoDocuments) {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task deleted successfully"})
}

func (tc *TaskController) Register(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdUser, err := tc.UserService.CreateUser(ctx, models.User{
		Username: payload.Username,
		Password: payload.Password,
		Role:     models.RoleUser,
	})
	if err != nil {
		switch {
		case errors.Is(err, data.ErrUserAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		case errors.Is(err, data.ErrInvalidUsername), errors.Is(err, data.ErrInvalidPassword):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "user registered successfully",
		"id":       createdUser.ID,
		"username": createdUser.Username,
		"role":     createdUser.Role,
	})
}

func (tc *TaskController) Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := tc.UserService.FindByUsername(ctx, loginData.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	token, err := middleware.GenerateToken(*user, 24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"token":   token,
		"role":    user.Role,
	})
}
