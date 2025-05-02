package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Team struct {
	bun.BaseModel `bun:"table:teams,alias:t"`

	ID          int64     `bun:"id,pk,autoincrement" json:"id"`
	Name        string    `bun:"name,notnull" json:"name"`
	Description string    `bun:"description" json:"description"`
	OwnerID     int64     `bun:"owner_id,notnull" json:"owner_id"`
	CreatedAt   time.Time `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt   time.Time `bun:"deleted_at,soft_delete" json:"-"`

	Owner   *User             `bun:"rel:belongs-to,join:owner_id=id" json:"owner,omitempty"`
	Members []*TeamMembership `bun:"rel:has-many,join:id=team_id" json:"members,omitempty"`
}

type TeamRole string

const (
	TeamRoleOwner  TeamRole = "owner"
	TeamRoleAdmin  TeamRole = "admin"
	TeamRoleMember TeamRole = "member"
)

type TeamMembership struct {
	bun.BaseModel `bun:"table:team_memberships,alias:tm"`

	ID        int64     `bun:"id,pk,autoincrement" json:"id"`
	TeamID    int64     `bun:"team_id,notnull" json:"team_id"`
	UserID    int64     `bun:"user_id,notnull" json:"user_id"`
	Role      TeamRole  `bun:"role,notnull" json:"role"`
	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt time.Time `bun:"deleted_at,soft_delete" json:"-"`

	Team *Team `bun:"rel:belongs-to,join:team_id=id" json:"team,omitempty"`
	User *User `bun:"rel:belongs-to,join:user_id=id" json:"user,omitempty"`
}
