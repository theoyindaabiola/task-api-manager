package middleware

import (
	"taskapi/dao"

	"github.com/gin-gonic/gin"
)

// This file is only for is only for permission checks
/*  
	Authenticate the permission to be created. 
	taskDAO *dao.TaskDAO gives access to the TaskDAO struct functions e.g GetTaskDB.
*/
func TaskOwnerMiddleware(taskDAO *dao.TaskDAO) gin.HandlerFunc {
    return func(c *gin.Context) { // all below form the context string
		// fetch and check task exists
        userID, exists := c.Get("user_id") // from JWTAuthMiddleware context 
        if !exists {
            c.JSON(401, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }

		// fetch and check task exists
        taskID := c.Param("id") // from JWTAuthMiddleware URL
        task, err := taskDAO.GetTaskDB(taskID)
        if err != nil {
            c.JSON(404, gin.H{"error": "Task not found"})
            c.Abort()
            return
        }

		// check user task is created by who is trying to access it
        if task.CreatedBy == userID {
            c.Next() // Owner
        } else {
			c.JSON(403, gin.H{"error": "Only task owner can delegate"})
			c.Abort()
		}
    }
}

/*  
	Authenticate the permission to be created. 
	permDAO *dao.TaskPermissionDAO gives access to the TaskPermissionDAO struct functions e.g FindPermission.
*/
func TaskAccessMiddleware(taskDAO *dao.TaskDAO, permDAO *dao.TaskPermissionDAO, permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr, exists := c.Get("user_id") // this returns any type
		if !exists {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// asserts the userID is a string for use in dao
		userID, ok := userIDStr.(string) // returns a string and boolean
        if !ok {
            c.JSON(500, gin.H{"error": "Invalid user ID type"})
            c.Abort()
            return
        }

		// fetch and checks if the task exists
		taskID := c.Param("id") // this returns a string
		task, err := taskDAO.GetTaskDB(taskID)
		if err != nil {
			c.JSON(404, gin.H{"error": "Task not found"})
			c.Abort()
			return
		}

		// allow access for the task owner
		if task.CreatedBy == userID {
			c.Next()
			return
		}

		// checking for permission
		perm, err := permDAO.FindPermission(taskID, userID) // this checks the permission table
		if err != nil || perm == nil {
			c.JSON(403, gin.H{"error": "No permission"})
			c.Abort()
			return
		}

		if (permission == "read" && perm.CanRead) || (permission == "update" && perm.CanUpdate){
			c.Next()
		} else {
			c.JSON(403, gin.H{"error": "Insufficient permissions"})
			c.Abort()
		}
	}
}