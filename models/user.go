package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// GORM
type User struct {
	// json is in the go application while the db is mapped to it in the database
	ID					uuid.UUID 		`gorm:"primarykey" json:"id"`
	Username 			string 			`gorm:"unique;not null" json:"username"`
	Password			string 			`gorm:"not null" json:"-"`
	Email				string 			`gorm:"unique;not null" json:"email"`
	PhoneNumber			*string			`gorm:"unique" json:"phone_number"` // optionally empty - nil == null
	Verified			bool   			`gorm:"default:false" json:"verified"`
	Enabled2FA			bool			`gorm:"default:false" json:"enabled_2fa"`				// "sms", "totp", "email"
	IsSmsVerified		bool			`gorm:"default:false" json:"is_sms_verified"`
	VerificationToken 	*string 		`gorm:"unique" json:"-"` 			// optionally empty
	ResetToken			*string			`gorm:"unique" json:"-"` 			// optionally empty
	SmsOTP				*string			`json:"-"`							// optionally empty
	SmsOTPExpiresAt		*time.Time		`gorm:"default:null" json:"-"`
	ResetTokenExpiresAt *time.Time      `gorm:"default:null" json:"-"`
	CreatedAt			time.Time 		`json:"created_at"`
	UpdatedAt   		time.Time 		`json:"updated_at"`
	DeletedAt   		gorm.DeletedAt 	`gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return nil
}
