package models

import (
	"database/sql"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Id string `gorm:"primaryKey;type:varchar(255)" json:"id"`

	// ouath info
	Email          string         `gorm:"uniqueIndex"`
	GoogleID       sql.NullString `gorm:"uniqueIndex;default:null"`
	DisplayName    string
	ProfilePicture string

	//zklogin fields
	Sub        sql.NullString `gorm:"uniqueIndex;default:null"`
	Salt       sql.NullString `gorm:"index;default:null"`
	SuiAddress sql.NullString `gorm:"uniqueIndex;default:null"`

	// manage sessions
	RefreshToken sql.NullString
	TokenExpiry  sql.NullTime

	// email and password fields
	PasswordHash  sql.NullString `gorm:"type:varchar(255)"`
	EmailVerified bool           `gorm:"default:false"`
	OTPCode       sql.NullString `gorm:"type:varchar(6)"`
	OTPExpiry     sql.NullTime
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.Id = uuid.NewString()

	return nil
}
