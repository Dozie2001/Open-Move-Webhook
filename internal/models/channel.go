package models

import (
	"time"

	"github.com/uptrace/bun"
)

type ChannelType string

const (
	ChannelTypeWebhook  ChannelType = "webhook"
	ChannelTypeEmail    ChannelType = "email"
	ChannelTypeTelegram ChannelType = "telegram"
	ChannelTypeDiscord  ChannelType = "discord"
)

type Channel struct {
	bun.BaseModel `bun:"table:channels,alias:c"`

	ID          int64       `bun:"id,pk,autoincrement" json:"id"`
	Name        string      `bun:"name,notnull" json:"name"`
	Description string      `bun:"description" json:"description"`
	Type        ChannelType `bun:"type,notnull" json:"type"`
	Config      string      `bun:"config,type:jsonb,notnull" json:"config"`
	TeamID      *int64      `bun:"team_id" json:"team_id,omitempty"`
	UserID      int64       `bun:"user_id,notnull" json:"user_id"`
	CreatedAt   time.Time   `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time   `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt   time.Time   `bun:"deleted_at,soft_delete" json:"-"`

	Team          *Team                  `bun:"rel:belongs-to,join:team_id=id" json:"team,omitempty"`
	User          *User                  `bun:"rel:belongs-to,join:user_id=id" json:"user,omitempty"`
	Subscriptions []*SubscriptionChannel `bun:"rel:has-many,join:id=channel_id" json:"subscriptions,omitempty"`
}

type SubscriptionChannel struct {
	bun.BaseModel `bun:"table:subscription_channels,alias:sc"`

	ID             int64     `bun:"id,pk,autoincrement" json:"id"`
	SubscriptionID int64     `bun:"subscription_id,notnull" json:"subscription_id"`
	ChannelID      int64     `bun:"channel_id,notnull" json:"channel_id"`
	CreatedAt      time.Time `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt      time.Time `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt      time.Time `bun:"deleted_at,soft_delete" json:"-"`

	Subscription *Subscription `bun:"rel:belongs-to,join:subscription_id=id" json:"subscription,omitempty"`
	Channel      *Channel      `bun:"rel:belongs-to,join:channel_id=id" json:"channel,omitempty"`
}
