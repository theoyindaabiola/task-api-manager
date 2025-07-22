package dao

import (
	"errors"
	"taskapi/models"

	"gorm.io/gorm"
)

// create an instance of the data access object
type TaskPermissionDAO struct {
	DB *gorm.DB
}

func NewTaskPermissionDAO(db *gorm.DB) *TaskPermissionDAO {
	// return the address of the instance of the db we are passing in
	return &TaskPermissionDAO{DB: db}
}

func (dao *TaskPermissionDAO) CreatePermission(taskPermission *models.TaskDelegation) error {
	// gorm needs the instance of Task{} not the task struct
	if err := dao.DB.Create(&taskPermission).Error; err != nil {
		return err
	}
	return nil
}

func (dao *TaskPermissionDAO) FindPermission(taskID string, userID string) (*models.TaskDelegation, error) {
	var taskPermission models.TaskDelegation
	err := dao.DB.Where("task_id = ? AND delegatee_id = ?", taskID, userID).First(&taskPermission).Error
	
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// return nil, nil when no record exists
		return nil, nil
	}
	
	if err != nil {
        return nil, err
    }

	return &taskPermission, nil
}

func HasPermission(required string, actual rune) bool {
    switch required {
    case "R":
        return actual == 'R' || actual == 'U' || actual == 'O'
    case "U":
        return actual == 'U' || actual == 'O'
    case "O":
        return actual == 'O'
    default:
        return false
    }
}

func (dao *TaskPermissionDAO) UpdatePermission(taskID, userID string, newPermission rune) error {
	return dao.DB.Model(&models.TaskDelegation{}).
	Where("task_id = ? AND delegatee_id = ?", taskID, userID).Update("permission", newPermission).Error
}

func (dao *TaskPermissionDAO) DeletePermission(taskID, userID string) error {
	err := dao.DB.Exec("DELETE FROM task_delegations WHERE task_id = ? AND delegatee_id = ?", taskID, userID).Error
	return err
}
