package router

import (
	"task_manager/controllers"
	"task_manager/data"

	"github.com/gin-gonic/gin"
)

func SetupRouter(service *data.TaskService) *gin.Engine {

	r := gin.Default()

	taskController := controllers.TaskController{
		Service: service,
	}

	r.GET("/tasks", taskController.GetTasks)
	r.GET("/tasks/:id", taskController.GetTask)
	r.POST("/tasks", taskController.CreateTask)
	r.PUT("/tasks/:id", taskController.UpdateTask)
	r.DELETE("/tasks/:id", taskController.DeleteTask)

	return r
}
