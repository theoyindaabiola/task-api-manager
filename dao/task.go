package dao

import (
	"taskapi/models"

	"gorm.io/gorm"
)

// create an instance of the data access object
type TaskDAO struct {
	DB *gorm.DB
}

// create the instance of the dao passing in the pointer
func NewTaskDAO(db *gorm.DB) *TaskDAO {
	// return the address of the instance of the db we are passing in
	return &TaskDAO{DB: db}
}

/** 
	to create task: pointer function: these functions are created 
	under the dao package of TaskDAO class, and the parameter is called dao
	These are methods under TaskDAO class e.g: TaskDAO.CreateTaskDB
	a pointer to the models is needed to connect db with the model structure, 
	because they are now on different packages
**/

func (dao *TaskDAO) CreateTaskDB(task *models.Task) error {
	// gorm needs the instance of Task{} not the task struct
	if err := dao.DB.Create(&task).Error; err != nil {
		return err
	}
	return nil
}

// to read task, get the array of the task
func (dao *TaskDAO) GetTasksDB() ([]models.Task, error) {
	var tasks [] models.Task
	if err := dao.DB.Find(&tasks).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}

// the ID string is coming from the API request URL
func (dao *TaskDAO) GetTaskDB(id string) (*models.Task, error) {
    var task models.Task
    if err := dao.DB.Where("id = ?", id).First(&task).Error; err != nil {
        return nil, err
    }
    return &task, nil
}

// here the GORM accepts map as struct for updating, empty interface is flexible for updating, struct is not.
func (dao *TaskDAO) UpdateTaskDB(taskID string, task map[string]interface{}) error {
	// placeholder for the task to be updated
	var updateTask models.Task
	// finds the task by id and store in the memory location of updateTask
	if err := dao.DB.Where("id = ?", taskID).First(&updateTask).Error; err != nil {
		return err
	}
	// return the fetched task to be updated and update the interface values of task
	return dao.DB.Model(&updateTask).Updates(task).Error

	// return dao.DB.Model(&models.Task{}).Where("id = ?", taskID).Updates(task).Error
}

func (dao *TaskDAO) DeleteTaskDB(id string) error {
	var task models.Task
	if err := dao.DB.Where("id = ?", id).First(&task).Error; err != nil {
		return err
	}

	return dao.DB.Delete(&task).Error
}
