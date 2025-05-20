package services

import (
	// "fmt"
	"os"
	"taskapi/dao"    // needs to interact with it
	"taskapi/models" // needs the model for the db
	"taskapi/utils"

	"github.com/google/uuid"
)

/**
	instance of the dao(database) named TaskService to be able to interact with the database
	so the functions has to be connected to this instance
**/

type TaskService struct {
	// package dao, struct TaskDAO, instance named TaskDAO
	TaskDAO 			*dao.TaskDAO
	TaskPermissionDAO 	*dao.TaskPermissionDAO // for permissions
    UserDAO         	*dao.UserDAO  
}

func NewTaskService(dao *dao.TaskDAO, taskPermissionDAO *dao.TaskPermissionDAO, userDAO *dao.UserDAO) *TaskService {
	// return an instance of the NewTaskDAO instance of the dao we are passing in
	return &TaskService{
		TaskDAO: 			dao,
		TaskPermissionDAO: 	taskPermissionDAO,
        UserDAO:         	userDAO,
	}
}

// the CreateTaskDB() needs the model task parameter
// pointing to the memory location of TaskService
func (s *TaskService) CreateTask(task *models.Task) error { 
	// CreateTaskDB is a function of class TaskDAO from the dao
	return s.TaskDAO.CreateTaskDB(task)
}

func (s *TaskService) GetTasks() ([]models.Task, error) { 
	// GetTasksDB is a function of class TaskDAO from the dao
	return s.TaskDAO.GetTasksDB()
}

// id string coming from the API request
func (s *TaskService) GetTask(id string) (*models.Task, error) { 
	// GetTaskDB is a function of class TaskDAO from the dao
	return s.TaskDAO.GetTaskDB(id)
}

func (s *TaskService) UpdateTask(taskID uuid.UUID, task map[string]interface{}) error { 
	// UpdateTaskDB is a function of class TaskDAO from the dao
	return s.TaskDAO.UpdateTaskDB(taskID, task)
}

// func (s *TaskService) UpdateTask(taskID string, task map[string]interface{}) error { 
// 	// UpdateTaskDB is a function of class TaskDAO from the dao
// 	return s.TaskDAO.UpdateTaskDB(taskID, task)
// }

func (s *TaskService) DelegateTask(permission *models.TaskPermission) error {
	if err := s.TaskPermissionDAO.CreatePermission(permission); err != nil {
        return err
    }

	// fetch task and user for notification
    task, err := s.TaskDAO.GetTaskDB(permission.TaskID)
    if err != nil {
        return  err
    }

	user, err := s.UserDAO.GetUserByIdDB(permission.UserID)
    if err != nil {
        return err
    }

	if err := utils.PublishMessage(
        os.Getenv("TASK_DELEGATION_QUEUE"),
        user.Email,
        permission.TaskID,
        "delegation",
        permission.UserID,
        user.Email,
        task.Title,
    ); err != nil {
        return err
    }
    return nil
}

func (s *TaskService) DeleteTask(id string) error { 
	// DeleteTaskDB is a function of class TaskDAO from the dao
	return s.TaskDAO.DeleteTaskDB(id)
}

