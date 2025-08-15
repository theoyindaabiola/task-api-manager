package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// GORM
type User struct {
	// json is in the go application while the db i s mapped to it in the database
	ID					uuid.UUID 		`gorm:"primarykey" json:"id"`
	Username 			string 			`gorm:"unique;not null" json:"username"`
	Password			string 			`gorm:"not null" json:"-"`
	Email				string 			`gorm:"unique;not null" json:"email"`
	Verified			bool   			`gorm:"default:false" json:"verified"`
	VerificationToken 	*string 		`gorm:"unique" json:"-"`
	ResetToken			*string			`gorm:"unique" json:"-"`
	ResetTokenExpiresAt *time.Time      `gorm:"default:null" json:"-"`
	CreatedAt			time.Time 		`json:"created_at"`
	UpdatedAt   		time.Time 		`json:"updated_at"`
	DeletedAt   		gorm.DeletedAt 	`gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return nil
}
