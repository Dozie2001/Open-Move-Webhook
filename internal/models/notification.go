package models

import (
	"time"

	"github.com/uptrace/bun"
)

type NotificationStatus string

const (
	NotificationStatusPending   NotificationStatus = "pending"
	NotificationStatusDelivered NotificationStatus = "delivered"
	NotificationStatusFailed    NotificationStatus = "failed"
)

type Notification struct {
	bun.BaseModel `bun:"table:notifications,alias:n"`

	ID               int64              `bun:"id,pk,autoincrement" json:"id"`
	SubscriptionID   int64              `bun:"subscription_id,notnull" json:"subscription_id"`
	ChannelID        int64              `bun:"channel_id,notnull" json:"channel_id"`
	EventPayload     string             `bun:"event_payload,type:jsonb,notnull" json:"event_payload"`
	DeliveryPayload  string             `bun:"delivery_payload,type:jsonb" json:"delivery_payload"`
	Status           NotificationStatus `bun:"status,notnull" json:"status"`
	ErrorMessage     string             `bun:"error_message" json:"error_message,omitempty"`
	DeliveryAttempts int                `bun:"delivery_attempts,notnull,default:0" json:"delivery_attempts"`
	CreatedAt        time.Time          `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt        time.Time          `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`
	DeliveredAt      *time.Time         `bun:"delivered_at" json:"delivered_at,omitempty"`

	Subscription *Subscription `bun:"rel:belongs-to,join:subscription_id=id" json:"subscription,omitempty"`
	Channel      *Channel      `bun:"rel:belongs-to,join:channel_id=id" json:"channel,omitempty"`
}
