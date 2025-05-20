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
	// group the routes
	taskRoutes := router.Group("/api/tasks")
	taskRoutes.Use(middleware.JWTAuthMiddleware()) // general middleware for authentication
	{
		taskRoutes.POST("/", taskController.CreateTask)
		// taskRoutes.GET("/", taskController.verify-email)
		taskRoutes.GET("/", taskController.GetTasks)
		// taskRoutes.GET("/:id", taskController.GetTask) // routes and http requests
		// taskRoutes.PUT("/:id", taskController.UpdateTask)
		taskRoutes.GET("/:id", middleware.TaskAccessMiddleware(taskDAO, taskPermissionDAO, "read"), taskController.GetTask) // middleware for authorization and permission
		taskRoutes.PUT("/:id", middleware.TaskAccessMiddleware(taskDAO, taskPermissionDAO, "update"), taskController.UpdateTask) // middleware for authorization and permission
		taskRoutes.DELETE("/:id", taskController.DeleteTask)
		taskRoutes.POST("/:id/delegate", middleware.TaskOwnerMiddleware(taskDAO), taskController.DelegateTask) // middleware for authorization
	}
}
// Compare this snippet from services/task.go:

// ??? Can't I use taskController all through? instead of the JWTAuthMiddleware, as a user how do I access the middleware directly?