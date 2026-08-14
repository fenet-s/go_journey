package controllers

import (
	"net/http"

	"refactored_task_management_api/domain"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ---- Task controller ----

type TaskController struct {
	taskUsecase domain.TaskUsecase
}

func NewTaskController(taskUsecase domain.TaskUsecase) *TaskController {
	return &TaskController{taskUsecase: taskUsecase}
}

func (tc *TaskController) CreateTask(c *gin.Context) {
	var task domain.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userHex := c.GetString("userID")
	userOID, err := primitive.ObjectIDFromHex(userHex)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in context"})
		return
	}
	task.UserID = userOID

	if err := tc.taskUsecase.CreateTask(c.Request.Context(), &task); err != nil {
		c.JSON(statusFor(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, task)
}

func (tc *TaskController) GetTask(c *gin.Context) {
	taskOID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID format"})
		return
	}

	task, err := tc.taskUsecase.GetTask(c.Request.Context(), taskOID, c.GetString("userID"), c.GetString("role"))
	if err != nil {
		c.JSON(statusFor(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (tc *TaskController) GetAllTasks(c *gin.Context) {
	tasks, err := tc.taskUsecase.GetAllTasks(c.Request.Context(), c.GetString("userID"), c.GetString("role"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (tc *TaskController) UpdateTask(c *gin.Context) {
	taskOID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID format"})
		return
	}

	var task domain.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task.ID = taskOID

	userHex := c.GetString("userID")
	if userOID, err := primitive.ObjectIDFromHex(userHex); err == nil {
		task.UserID = userOID
	}

	if err := tc.taskUsecase.UpdateTask(c.Request.Context(), &task, c.GetString("userID"), c.GetString("role")); err != nil {
		c.JSON(statusFor(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (tc *TaskController) DeleteTask(c *gin.Context) {
	taskOID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID format"})
		return
	}

	if err := tc.taskUsecase.DeleteTask(c.Request.Context(), taskOID, c.GetString("userID"), c.GetString("role")); err != nil {
		c.JSON(statusFor(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "task deleted"})
}

// ---- User controller ----

type UserController struct {
	userUsecase domain.UserUsecase
}

func NewUserController(userUsecase domain.UserUsecase) *UserController {
	return &UserController{userUsecase: userUsecase}
}

func (uc *UserController) Register(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uc.userUsecase.Register(c.Request.Context(), &user); err != nil {
		c.JSON(statusFor(err), gin.H{"error": err.Error()})
		return
	}

	type response struct {
		Message string `json:"message"`
		ID      string `json:"id"`
		Email   string `json:"email"`
		Role    string `json:"role"`
	}
	c.JSON(http.StatusCreated, response{
		Message: "user registered successfully",
		ID:      user.ID.Hex(),
		Email:   user.Email,
		Role:    user.Role,
	})
}

func (uc *UserController) Login(c *gin.Context) {
	var creds struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := uc.userUsecase.Login(c.Request.Context(), creds.Email, creds.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// statusFor maps domain errors to HTTP status codes
func statusFor(err error) int {
	switch err {
	case domain.ErrTaskNotFound, domain.ErrUserNotFound:
		return http.StatusNotFound
	case domain.ErrUnauthorized:
		return http.StatusForbidden
	case domain.ErrUserExists, domain.ErrEmptyTitle, domain.ErrInvalidCreds:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
