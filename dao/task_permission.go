package dao

import (
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
	if err := dao.DB.Where("task_id = ? AND delegatee_id = ?", taskID, userID).First(&taskPermission).Error; err != nil {
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


