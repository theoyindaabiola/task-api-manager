package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// this links user with the permitted tasks
type TaskDelegation struct {
	TaskID			string 		`gorm:"primaryKey" json:"task_id"`
	DelegateeID		string 		`gorm:"primaryKey" json:"delegatee_id"`
    Permission		rune		`gorm:"not null" json:"permission"`
	Task			Task		`gorm:"foreignKey:TaskID" json:"task"`
	Delegatee		User  		`gorm:"foreignKey:DelegateeID;" json:"delegatee"`
	CreatedAt		time.Time 	`json:"created_at"`
	UpdatedAt		time.Time 	`json:"updated_at"`
}

func (tp *TaskDelegation) BeforeCreate(tx *gorm.DB) (err error) {
    if tp.TaskID == "" || tp.DelegateeID == "" {
        return errors.New("TaskID and UserID must be set")
    }
    return nil
}
