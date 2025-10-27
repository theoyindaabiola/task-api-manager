package routes

import (
	"taskapi/controllers"
	"taskapi/dao"
	"taskapi/middleware"
	"taskapi/services"

	"github.com/gin-gonic/gin"
)

type Routes struct {
	TaskController *controllers.TaskController
}

func RegisterRoutes(router *gin.Engine, taskController *controllers.TaskController, taskDAO *dao.TaskDAO, taskPermissionDAO *dao.TaskPermissionDAO, userService *services.UserService) {
	taskRoutes := router.Group("/api/tasks")
	taskRoutes.Use(middleware.JWTAuthMiddleware(userService))
	{
		taskRoutes.POST("/", taskController.CreateTask)
		taskRoutes.GET("/", taskController.GetTasks)

		taskRoutes.GET("/:id", middleware.TaskAccessMiddleware(taskDAO, taskPermissionDAO, "R"), taskController.GetTask)
		taskRoutes.PUT("/:id", middleware.TaskAccessMiddleware(taskDAO, taskPermissionDAO, "U"), taskController.UpdateTask)

		taskRoutes.POST("/:id/delegate", middleware.TaskOwnerMiddleware(taskDAO, taskPermissionDAO, "O"), taskController.DelegateTask)
		taskRoutes.DELETE("/:id/delete", middleware.TaskOwnerMiddleware(taskDAO, taskPermissionDAO, "O"), taskController.DeleteTask)

		taskRoutes.PATCH("/:id/permission", middleware.TaskOwnerMiddleware(taskDAO, taskPermissionDAO, "O"), taskController.UpdatePermission)
		taskRoutes.DELETE("/:id/permission", middleware.TaskOwnerMiddleware(taskDAO, taskPermissionDAO, "O"), taskController.RevokePermission)

	}
}
