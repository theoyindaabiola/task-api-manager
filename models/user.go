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
	Verified			bool   			`gorm:"default:false" json:"verified"`
	Enabled2FA			bool			`gorm:"default:false" json:"enabled_2fa"`
	IsOtpVerified		bool			`gorm:"default:false" json:"is-otp-verified"`
	IsTotpVerified		bool			`gorm:"default:false" json:"is_totp_verified"`
	VerificationToken 	*string 		`gorm:"unique" json:"-"`
	ResetToken			*string			`gorm:"unique" json:"-"`
	TOTPSecret			*string			`json:"_"` 		// store secret for authenticator apps
	ResetTokenExpiresAt *time.Time      `gorm:"default:null" json:"-"`
	CreatedAt			time.Time 		`json:"created_at"`
	UpdatedAt   		time.Time 		`json:"updated_at"`
	DeletedAt   		gorm.DeletedAt 	`gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return nil
}
