package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GORM
type Task struct {
	// json is in the go application while the db i s mapped to it in the database
	// ID			uuid.UUID `gorm:"type:uuid;default:uuid_default_generate_v4();"primaryKey" json:"id"` 
	ID			string 		`gorm:"primaryKey" json:"id"`
	Title 		string 		`gorm:"not null" json:"title"`
	Description	string 		`gorm:"not null" json:"description"`
	Completed 	bool 		`gorm:"default:false" json:"completed"`
	CreatedBy	string 		`json:"created_by"`
	CreatedAt 	time.Time 	`json:"created_at"`
	UpdatedAt 	time.Time 	`json:"updated_at"`
}

func (t *Task) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.NewString()
	}	
	return nil
}

