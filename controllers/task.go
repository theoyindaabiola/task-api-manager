package controllers

import (
	"net/http"
	"taskapi/dto"
	"taskapi/models"
	"taskapi/services"

	"github.com/gin-gonic/gin"
)

/**
	service instance needed to send the payload to it and the functions
	has to be connected to this instance
**/

type TaskController struct {
	TaskService *services.TaskService
	UserService *services.UserService
}

func NewTaskController(taskService *services.TaskService, userService *services.UserService) *TaskController {
	return &TaskController{TaskService: taskService, UserService: userService}
}

// context gets into the body of your package. c is an instance of pointing to the gin.Context
func (tc *TaskController) CreateTask(c *gin.Context) {
	// security: this avoid users ability to override models.task parameters
	var taskInput struct {
        Title       string `json:"title" binding:"required"`
        Description string `json:"description"`
    }
	
	// read the request JSON data and convert it to task of struct Task
	if err := c.ShouldBindJSON(&taskInput); err != nil {
		// if error, return error using http in JSON format using the gin context
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) 
		return
	}

	// get user_id from context
    userID, _ := c.Get("user_id")

	// asserts the userID is a string
	userIDStr, ok := userID.(string) 
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
        return
    }

	// create task
    task := models.Task{
        Title:       taskInput.Title,
        Description: taskInput.Description,
        CreatedBy:   userIDStr,
    }

	/**
		creates the database and return error in JSON format
		tc.TaskService connects to the services and call the CreateTask()
	**/
	if err := tc.TaskService.CreateTask(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, "Task successfully created")
}

// c is an instance of gin.Context, and the instance is needed
func (tc *TaskController) GetTasks(c *gin.Context) {
	// call the services then database
	tasks, err := tc.TaskService.GetTasks()

	// return error if there is an issue with getting response from the database
	if err != nil { 
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// gin handles everything about http request and response, it is connected to HTTP
func (tc *TaskController) GetTask(c *gin.Context) { 
	taskID := c.Param("id")

	userID, _ := c.Get("user_id")
	// asserts the userID is a string
	userIDStr, ok := userID.(string) 
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
        return
    }

	task, err := tc.TaskService.GetTask(taskID, userIDStr)
	// return error message if task not found
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Task not found"})
	}
	c.JSON(http.StatusOK, task)
}

func (tc *TaskController) UpdateTask(c *gin.Context) {
	// needed as the key, coming from the URL request
	taskID := c.Param("id")

	userID, _ := c.Get("user_id")
	userIDStr, ok := userID.(string)
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
        return
    }

	// payload placeholder
	var task map[string]interface{}
	// get and confirm that there is no error with the payload
	if err := c.ShouldBindJSON(&task); err != nil {
		// if error, return error using http in JSON format using the gin context
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) 
		return
	}

	if err := tc.TaskService.UpdateTask(taskID, userIDStr, task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

/** 
	DelegateTask is a complex business logic
	Requires multiple validations across the daos(user, task & permission)
**/ 
func (tc *TaskController) DelegateTask(c *gin.Context) {
	taskID := c.Param("id")

	var delegateInput dto.TaskDelegationInput

	// bind JSON body to input struct (i.e of delegated/user's id)
	if err := c.ShouldBindJSON(&delegateInput); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	userIDStr, ok := userID.(string)
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
        return
    }

	// verify task exists
	_, err := tc.TaskService.GetTask(taskID, userIDStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// verify delegated user exists
	if _, err := tc.UserService.GetUserByID(delegateInput.DelegateeID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
        return
    }

	// validate permission input (must be 1 character)
	if len(delegateInput.Permission) != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Permission must be a single character 'R' or 'U'"})
		return
	}

	permRune := rune(delegateInput.Permission[0])

	// set and create TaskPermission entry for the delegated user
	permission := &models.TaskDelegation {
		TaskID:    		taskID,
		DelegateeID:    delegateInput.DelegateeID,
		Permission:  	permRune,
	}
	if err := tc.TaskService.DelegateTask(permission); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Task delegated successfully"})
}

func (tc *TaskController) UpdatePermission(c *gin.Context) {
	taskID := c.Param("id")

	var input dto.TaskDelegationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	if len(input.Permission) != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Permission must be a single character like 'R' or 'U'"})
		return
	}
	permRune := rune(input.Permission[0])

	// get user ID of the owner (JWT context)
	userID, _ := c.Get("user_id")
	ownerID, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	// only task owner can update permissions
	_, err := tc.TaskService.GetTask(taskID, ownerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// call the service method to update
	if err := tc.TaskService.UpdateTaskPermission(taskID, input.DelegateeID, permRune); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission updated successfully"})
}

func (tc *TaskController) RevokePermission(c *gin.Context) {
	taskID := c.Param("id")
	userID, _ := c.Get("user_id")
	ownerID, ok := userID.(string)
	if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
        return
    }

	var body struct {
		DelegateeID string `json:"delegatee_id"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := tc.TaskService.RevokePermission(taskID, ownerID, body.DelegateeID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission revoked successfully"})
}

func (tc *TaskController) DeleteTask(c *gin.Context) {
	taskID := c.Param("id")

	userID, _ := c.Get("user_id")
	userIDStr, ok := userID.(string) // asserts the userID is a string
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
        return
    }

	if err := tc.TaskService.DeleteTask(taskID, userIDStr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, "Task deleted succefully")
}
