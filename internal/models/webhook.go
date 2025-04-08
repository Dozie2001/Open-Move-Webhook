package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Webhook struct {
	Id        string    `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Url       string    `gorm:"column:url" json:"url"`
	Secret    string    `gorm:"column:secret" json:"secret"`
	Status    bool      `json:"status" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (w *Webhook) BeforeCreate(tx *gorm.DB) error {
	w.Id = uuid.NewString()

	return nil
}
