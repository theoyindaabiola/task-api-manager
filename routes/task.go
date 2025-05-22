package routes

import (
	"taskapi/dao"
	"taskapi/controllers"
	"taskapi/middleware"
	"github.com/gin-gonic/gin"
)

type Routes struct {
	TaskController *controllers.TaskController
}

func RegisterRoutes(router *gin.Engine, taskController *controllers.TaskController, taskDAO *dao.TaskDAO, taskPermissionDAO *dao.TaskPermissionDAO) {
	taskRoutes := router.Group("/api/tasks")
	taskRoutes.Use(middleware.JWTAuthMiddleware())
	{
		taskRoutes.POST("/", taskController.CreateTask)
		taskRoutes.GET("/", taskController.GetTasks)

		taskRoutes.GET("/:id", middleware.TaskAccessMiddleware(taskDAO, taskPermissionDAO, "R"), taskController.GetTask)
		taskRoutes.PUT("/:id", middleware.TaskAccessMiddleware(taskDAO, taskPermissionDAO, "U"), taskController.UpdateTask)

		taskRoutes.POST("/:id/delegate", middleware.TaskOwnerMiddleware(taskDAO, taskPermissionDAO, "O"), taskController.DelegateTask)
		taskRoutes.DELETE("/:id", middleware.TaskOwnerMiddleware(taskDAO, taskPermissionDAO, "O"), taskController.DeleteTask)
	}
}
// Compare this snippet from services/task.go:

// ??? Can't I use taskController all through? instead of the JWTAuthMiddleware, as a user how do I access the middleware directly?