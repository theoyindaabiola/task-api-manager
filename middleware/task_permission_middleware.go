package middleware

import (
	"taskapi/dao"

	"github.com/gin-gonic/gin"
)

/*
	This file is only for is only for permission checks
	Authenticate the permission to be created.
	taskDAO *dao.TaskDAO gives access to the TaskDAO struct functions e.g GetTaskDB.
*/
func TaskOwnerMiddleware(taskDAO *dao.TaskDAO, permDAO *dao.TaskPermissionDAO, permission string) gin.HandlerFunc {
    return func(c *gin.Context) { // all below form the context string
		// fetch and check task exists
		taskID := c.Param("id") // from JWTAuthMiddleware URL
		task, err := taskDAO.GetTaskDB(taskID)
		if err != nil {
			c.JSON(404, gin.H{"error": "Task not found"})
			c.Abort()
			return
		}

		// fetch and check user exists
		userIDStr, exists := c.Get("user_id")
		if !exists {
			c.JSON(404, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		userID, ok := userIDStr.(string)
		if !ok {
			c.JSON(500, gin.H{"error": "Invalid user ID in context"})
			c.Abort()
			return
		}

		// check user task is created by who is trying to access it
		if task.CreatedBy == userID {
			// task owner, assign 'O' permission
			if dao.HasPermission(permission, 'O') {
				c.Next() // owner is passed
				return
			} 
		} else {
			switch c.Request.Method {
			case "DELETE":
				c.JSON(403, gin.H{"error": "Only task owner can delete task"})
			case "PATCH":
				c.JSON(403, gin.H{"error": "Only task owner can update permissions"})
			default:
				c.JSON(403, gin.H{"error": "Only task owner can delegate"})
			}
			c.Abort() // non-owner is unauthorize
			return
		}
    }
}

/*  
	Authenticate the permission to be created. 
	permDAO *dao.TaskPermissionDAO gives access to the TaskPermissionDAO struct functions e.g FindPermission.
*/
func TaskAccessMiddleware(taskDAO *dao.TaskDAO, permDAO *dao.TaskPermissionDAO, permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// fetch and checks if the task exists
		taskID := c.Param("id") // this returns a string
		task, err := taskDAO.GetTaskDB(taskID)
		if err != nil {
			c.JSON(404, gin.H{"error": "Task not found"})
			c.Abort()
			return
		}

		// fetch and check user exists
		userIDStr, exists := c.Get("user_id")
		if !exists {
			c.JSON(404, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		userID, ok := userIDStr.(string)
		if !ok {
			c.JSON(500, gin.H{"error": "Invalid user ID in context"})
			c.Abort()
			return
		}

		// check user task is created by who is trying to access it
		if task.CreatedBy == userID {
			if dao.HasPermission(permission, 'O') {
				c.Next()
				return
			}
		}

		// Check if the user has been delegated permission
		perm, err := permDAO.FindPermission(taskID, userID)
		if err != nil {
			c.JSON(500, gin.H{"error": "Error checking permissions"})
			c.Abort()
			return
		}

		if perm == nil {
			c.JSON(403, gin.H{"error": "No task delegation or permission found"})
			c.Abort()
			return
		}

		if !dao.HasPermission(permission, perm.Permission) {
			c.JSON(403, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
