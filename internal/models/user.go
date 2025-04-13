package models

import (
	"database/sql"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Id string `gorm:"primaryKey;type:varchar(255)" json:"id"`

	// ouath info
	Email          string         `gorm:"uniqueIndex" json:"email"`
	GoogleID       sql.NullString `gorm:"uniqueIndex;default:null" json:"google_id"`
	DisplayName    string         `json:"display_name"`
	ProfilePicture string         `json:"profile_picture"`

	//zklogin fields
	Sub        sql.NullString `gorm:"uniqueIndex;default:null" json:"sub"`
	Salt       sql.NullString `gorm:"index;default:null" json:"salt"`
	SuiAddress sql.NullString `gorm:"uniqueIndex;default:null" json:"sui_address"`

	// manage sessions
	RefreshToken sql.NullString `json:"refresh_token"`
	TokenExpiry  sql.NullTime   `json:"token_expiry"`

	// email and password fields
	PasswordHash  sql.NullString `gorm:"type:varchar(255)" json:"password_hash"`
	EmailVerified bool           `gorm:"default:false" json:"email_verified"`
	OTPCode       sql.NullString `gorm:"type:varchar(6)" json:"otp_code"`
	OTPExpiry     sql.NullTime   `json:"otp_expiry"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.Id = uuid.NewString()

	return nil
}
