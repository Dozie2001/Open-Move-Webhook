package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Subscription struct {
	bun.BaseModel `bun:"table:subscriptions,alias:s"`

	ID          int64     `bun:"id,pk,autoincrement" json:"id"`
	Name        string    `bun:"name,notnull" json:"name"`
	Description string    `bun:"description" json:"description"`
	EventType   string    `bun:"event_type,notnull" json:"event_type"`
	TeamID      *int64    `bun:"team_id" json:"team_id,omitempty"`
	UserID      int64     `bun:"user_id,notnull" json:"user_id"`
	IsActive    bool      `bun:"is_active,notnull,default:true" json:"is_active"`
	CreatedAt   time.Time `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt   time.Time `bun:"deleted_at,soft_delete" json:"-"`

	Team     *Team                  `bun:"rel:belongs-to,join:team_id=id" json:"team,omitempty"`
	User     *User                  `bun:"rel:belongs-to,join:user_id=id" json:"user,omitempty"`
	Channels []*SubscriptionChannel `bun:"rel:has-many,join:id=subscription_id" json:"channels,omitempty"`
}
