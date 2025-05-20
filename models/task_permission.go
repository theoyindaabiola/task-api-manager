package models

import (
	"errors"
	"time"

	// "github.com/google/uuid"
	"gorm.io/gorm"
)

// this links user with the permitted tasks
type TaskPermission struct {
	TaskID			string 		`gorm:"primaryKey" json:"task_id"`
	UserID			string 		`gorm:"primaryKey" json:"user_id"`
	CanRead			bool   		`gorm:"not null" json:"can_read"`
	CanUpdate		bool   		`gorm:"not null" json:"can_update"`
	CreatedAt		time.Time 	`json:"created_at"`
	UpdatedAt		time.Time 	`json:"updated_at"`
}

func (tp *TaskPermission) BeforeCreate(tx *gorm.DB) (err error) {
    if tp.TaskID == "" || tp.UserID == "" {
        return errors.New("TaskID and UserID must be set")
    }
    return nil
}
