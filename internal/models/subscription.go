package models

import (
	// "encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type NewSubscriptionReqBody struct {
	WebhookId      string         `json:"webhook_id" binding:"required"`
	EventType      string         `json:"event_type" binding:"required"`
	FilterCriteria datatypes.JSON `json:"filter_criteria" binding:"required"`
}

type Subscription struct {
	Id             string         `gorm:"primaryKey;type:varchar(255)" json:"id"`
	WebhookId      string         `json:"webhook_id" gorm:"type:varchar(255)"`
	EventType      string         `json:"event_type"`
	FilterCriteria datatypes.JSON `json:"filter_criteria"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func (s *Subscription) BeforeCreate(tx *gorm.DB) error {
	s.Id = uuid.NewString()

	return nil
}
