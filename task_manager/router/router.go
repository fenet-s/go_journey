package router

import (
	"task_manager/controllers"
	"task_manager/data"
	"task_manager/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(service *data.TaskService, userService *data.UserService) *gin.Engine {

	r := gin.Default()

	taskController := controllers.TaskController{
		Service:     service,
		UserService: userService,
	}
	r.POST("/login", taskController.Login)
	r.POST("/register", taskController.Register)
	protected := r.Group("/tasks")
	protected.Use(middleware.AuthMiddleware())

	protected.GET("", taskController.GetTasks)
	protected.GET("/:id", taskController.GetTask)
	protected.POST("", taskController.CreateTask)
	protected.PUT("/:id", taskController.UpdateTask)
	admin := protected.Group("")
	admin.Use(middleware.RequireRole("admin"))
	admin.DELETE("/:id", taskController.DeleteTask)

	return r
}
