package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeliveryLog struct {
	Id           string    `gorm:"primaryKey;type:varchar(255)" json:"id"`
	WebhookId    string    `json:"webhook_id" gorm:"type:varchar(255)"`
	EventId      string    `json:"event_id" gorm:"type:varchar(255)"`
	Status       string    `json:"status" gorm:"default:false"`
	Attempts     int       `json:"attempts" gorm:"default:0"`
	LastAttempt  time.Time `json:"last_attempt"`
	NextRetry    time.Time `json:"next_retry"`
	ResponseCode int       `json:"response_code" gorm:"default:0"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (d *DeliveryLog) BeforeCreate(tx *gorm.DB) error {
	d.Id = uuid.NewString()

	return nil
}
