package routers

import (
	"refactored_task_management_api/delivery/controllers"
	"refactored_task_management_api/domain"
	"refactored_task_management_api/infrastructure"

	"github.com/gin-gonic/gin"
)

func SetupRouter(taskController *controllers.TaskController, userController *controllers.UserController, jwtService domain.JWTService) *gin.Engine {
	router := gin.Default()

	public := router.Group("/")
	{
		public.POST("/register", userController.Register)
		public.POST("/login", userController.Login)
	}

	protected := router.Group("/")
	protected.Use(infrastructure.AuthMiddleware(jwtService))
	{
		protected.POST("/tasks", taskController.CreateTask)
		protected.GET("/tasks", taskController.GetAllTasks)
		protected.GET("/tasks/:id", taskController.GetTask)
		protected.PUT("/tasks/:id", taskController.UpdateTask)
		admin := protected.Group("")
		admin.Use(infrastructure.RequireRole("admin"))
		admin.DELETE("/tasks/:id", taskController.DeleteTask)
	}

	return router
}
