package services

import (
	// "fmt"
	"errors"
	"os"
	"taskapi/dao"    // needs to interact with it
	"taskapi/models" // needs the model for the db
	"taskapi/utils"

	// "github.com/google/uuid"
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
func (s *TaskService) GetTask(taskID, userID string) (*models.Task, error) { 
	if err := s.AskPermission(taskID, userID); err != nil {
		return nil, err
	}
	// GetTaskDB is a function of class TaskDAO from the dao
	return s.TaskDAO.GetTaskDB(taskID)
}

func (s *TaskService) UpdateTask(taskID, userID string, task map[string]interface{}) error { 
	if err := s.AskPermission(taskID, userID); err != nil {
		return err
	}
	// UpdateTaskDB is a function of class TaskDAO from the dao if permission passed
	return s.TaskDAO.UpdateTaskDB(taskID, task)
}

func (s *TaskService) AskPermission(taskID string, userID string) error { 
	// check for task
	task, err := s.TaskDAO.GetTaskDB(taskID)
    if err != nil {
        return  err
    }

	if task.CreatedBy == userID {
		return nil // don't return any error
	}

	permission, err := s.TaskPermissionDAO.FindPermission(taskID, userID)
	if err != nil {
		return errors.New("no permission found for this task")
	}

	if permission.Permission != 'R' && permission.Permission != 'U' {
		return errors.New("insufficient permission")
	}
	return nil
}

func (s *TaskService) DelegateTask(permission *models.TaskDelegation) error {
	if err := s.TaskPermissionDAO.CreatePermission(permission); err != nil {
        return err
    }

	// fetch task and user for notification
    task, err := s.TaskDAO.GetTaskDB(permission.TaskID)
    if err != nil {
        return  err
    }

	user, err := s.UserDAO.GetUserByIdDB(permission.DelegateeID)
    if err != nil {
        return err
    }

	if err := utils.PublishMessage(
        os.Getenv("TASK_DELEGATION_QUEUE"),
        user.Email,
        permission.TaskID,
        "delegation",
        permission.DelegateeID,
        user.Email,
        task.Title,
    ); err != nil {
        return err
    }
    return nil
}

func (s *TaskService) DeleteTask(taskID, userID string) error {
	if err := s.AskPermission(taskID, userID); err != nil {
		return err
	}
	// DeleteTaskDB is a function of class TaskDAO from the dao
	return s.TaskDAO.DeleteTaskDB(taskID)
}

