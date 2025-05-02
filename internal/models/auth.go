package models

import (
	"time"

	"github.com/uptrace/bun"
)

type PasswordReset struct {
	bun.BaseModel `bun:"table:password_resets,alias:pr"`

	ID        int64     `bun:"id,pk,autoincrement" json:"-"`
	UserID    int64     `bun:"user_id,notnull" json:"-"`
	Token     string    `bun:"token,notnull,unique" json:"token"`
	ExpiresAt time.Time `bun:"expires_at,notnull" json:"expires_at"`
	Used      bool      `bun:"used,notnull,default:false" json:"-"`
	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp" json:"-"`

	User *User `bun:"rel:belongs-to,join:user_id=id" json:"-"`
}

type EmailVerification struct {
	bun.BaseModel `bun:"table:email_verifications,alias:ev"`

	ID        int64     `bun:"id,pk,autoincrement" json:"-"`
	UserID    int64     `bun:"user_id,notnull" json:"-"`
	Token     string    `bun:"token,notnull,unique" json:"token"`
	ExpiresAt time.Time `bun:"expires_at,notnull" json:"expires_at"`
	Used      bool      `bun:"used,notnull,default:false" json:"-"`
	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp" json:"-"`

	User *User `bun:"rel:belongs-to,join:user_id=id" json:"-"`
}
