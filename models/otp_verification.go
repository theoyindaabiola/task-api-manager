package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)


type OtpVerification struct {
	ID			uuid.UUID		`gorm:"primarykey" json:"id"`
	UserID		uuid.UUID 		`gorm:"index;not null" json:"user_id"`
	OtpCode		string 			`gorm:"not null" json:"-"` // multi-users can have same otp, so not unique
	ExpiresAt	time.Time 		`json:"expired_at"`
	OtpVerified	bool			`gorm:"default:false" json:"verified"`
	CreatedAt	time.Time 		`json:"created_at"`
}

func (u *OtpVerification) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return nil
}
